package middleware


import (
    "fmt"
    "os"
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
    tokens, ok := c.GetReqHeaders()["X-Api-Token"]
    if !ok || len(tokens) == 0 {
        fmt.Println("token not present in header")
        return fmt.Errorf("unauthorised")
    }

    token := tokens[0]
    
    claims, err := validateToken(token)
    if err != nil {
        return err
    }
    
    expiresFloat := claims["expires"].(float64)
    expires := int64(expiresFloat)
    //check token expiration
    if time.Now().Unix() > expires {
        return fmt.Errorf("token expired")
    }
   
    return c.Next()
}


func validateToken(tokenStr string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        fmt.Println("invalid signing method", token.Header["alg"])
        return nil, fmt.Errorf("unauthorised")
}
    secret := os.Getenv("JWT_SECRET")
    return []byte(secret), nil
    })

    if err != nil {
        fmt.Println("failed to parse JWT token:", err)
        return nil, fmt.Errorf("unwauthorised")
    }

    if !token.Valid {
        fmt.Println("invalid token")
        return nil, fmt.Errorf("unwauthorised")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, fmt.Errorf("unauthorised")
    }

    return claims, nil
}
