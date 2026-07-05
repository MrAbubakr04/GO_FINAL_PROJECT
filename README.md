Логика слоёв (очень важно)
1. domain (ядро, без зависимостей)

Тут только бизнес-сущности:

type Account struct {
    ID      int64
    UserID  int64
    Balance float64
}

❌ Никакого SQL, HTTP, JSON бизнес-логики
2. usecase (бизнес-логика)

Пример:

type AccountService struct {
    repo AccountRepository
}

func (s *AccountService) Deposit(userID int64, amount float64) error {
    if amount <= 0 {
        return errors.New("invalid amount")
    }

    return s.repo.AddBalance(userID, amount)
}

👉 Здесь вся логика кошелька:

перевод
пополнение
проверка баланса
3. repository (работа с БД)

Интерфейс:

type AccountRepository interface {
    GetByUserID(userID int64) (*domain.Account, error)
    AddBalance(userID int64, amount float64) error
}

Реализация PostgreSQL:

type accountRepo struct {
    db *sql.DB
}
4. delivery (HTTP слой)

Handler:

func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
    // parse request
    // call usecase
    // return response
}
5. cmd/api/main.go (точка входа)
func main() {
    db := postgres.Connect()

    accountRepo := repository.NewAccountRepo(db)
    accountService := usecase.NewAccountService(accountRepo)

    handler := handler.NewAccountHandler(accountService)

    router := http.NewRouter(handler)

    http.ListenAndServe(":8080", router)
}
🔥 Поток запроса (важно понимать)
HTTP Request
   ↓
delivery (handler)
   ↓
usecase (business logic)
   ↓
repository (DB)
   ↓
PostgreSQL
💡 Минимальный старт (если хочешь проще)

Если сейчас тяжело — начни так:

cmd/
internal/
    domain/
    usecase/
    repository/
    delivery/

И только потом дроби дальше.