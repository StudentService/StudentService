package application

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"backend/internal/domain/auth"
	"backend/internal/domain/user"
)

type AuthService struct {
	userRepo        user.Repository
	jwtSecret       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthService(userRepo user.Repository) *AuthService {
	// В продакшене берем из env
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-here-change-in-production"
	}

	accessTTL, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TTL"))
	if accessTTL == 0 {
		accessTTL = 15 // 15 минут по умолчанию
	}

	refreshTTL, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TTL"))
	if refreshTTL == 0 {
		refreshTTL = 24 * 7 // 7 дней по умолчанию
	}

	return &AuthService{
		userRepo:        userRepo,
		jwtSecret:       []byte(secret),
		accessTokenTTL:  time.Duration(accessTTL) * time.Minute,
		refreshTokenTTL: time.Duration(refreshTTL) * time.Hour,
	}
}

func (s *AuthService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.AuthResponse, error) {
	// 1. Находим пользователя по email
	u, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if u == nil {
		return nil, errors.New("invalid credentials")
	}

	// 2. Проверяем пароль с помощью метода модели
	if !u.CheckPassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// 3. Генерируем токены
	return s.generateTokenPair(u)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*auth.AuthResponse, error) {
	// 1. Валидируем refresh token
	claims := &auth.Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// 2. Проверяем, что это действительно refresh token (можно добавить поле type в claims)
	// 3. Получаем пользователя
	u, err := s.userRepo.GetByID(ctx, claims.UserID.String())
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 4. Генерируем новую пару токенов
	return s.generateTokenPair(u)
}

func (s *AuthService) generateTokenPair(u *user.User) (*auth.AuthResponse, error) {
	// Access token
	accessClaims := &auth.Claims{
		UserID: u.ID,
		Role:   string(u.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   u.ID.String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshClaims := &auth.Claims{
		UserID: u.ID,
		Role:   string(u.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   u.ID.String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &auth.AuthResponse{
		AccessToken:  accessString,
		RefreshToken: refreshString,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.accessTokenTTL.Seconds()),
	}, nil
}

// HashPassword - вспомогательная функция для хеширования паролей
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	// 1. Проверяем, что пользователь с таким email не существует
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("failed to check existing user")
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// 2. Проверяем, что пользователь с таким username не существует
	existingUser, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("failed to check existing username")
	}
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// 3. Создаем нового пользователя
	// Определяем роль (по умолчанию - student)
	role := user.RoleStudent
	if req.Role != "" {
		role = user.Role(req.Role)
	}

	// Создаем объект пользователя
	newUser := &user.User{
		Username:  req.Username,
		Email:     req.Email,
		Role:      role,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// 4. Устанавливаем пароль (хешируется внутри метода)
	if err := newUser.SetPassword(req.Password); err != nil {
		return nil, errors.New("failed to hash password: " + err.Error())
	}

	// 5. Сохраняем в БД
	createdUser, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, errors.New("failed to create user: " + err.Error())
	}

	// 6. Генерируем токены для автоматического входа
	tokens, err := s.generateTokenPair(createdUser)
	if err != nil {
		// В production лучше залогировать ошибку
		// Но пользователю возвращаем только данные пользователя
		return &auth.RegisterResponse{
			User:  createdUser.ToResponse(),
			Token: nil,
		}, nil // Не возвращаем ошибку, т.к. пользователь уже создан
	}

	// 7. Возвращаем успешный ответ с токенами
	return &auth.RegisterResponse{
		User:  createdUser.ToResponse(),
		Token: tokens,
	}, nil
}
