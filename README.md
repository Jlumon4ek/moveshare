
# Генерируем приватный ключ (2048 бит)
openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048

# Получаем публичный ключ из приватного
openssl rsa -pubout -in private.pem -out public.pem

# Генерация документации 
swag init -g cmd/server/main.go --parseDependency --parseInternal 