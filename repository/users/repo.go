package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/Zhanbatyr06/PizzaGoland/models"
	"github.com/Zhanbatyr06/PizzaGoland/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type Repo struct {
	coll *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{db.Collection("users")}
}

func (r Repo) Create(user *models.User) error {

	utils.Logger.WithField("action", "add_user").Info("Adding new user")
	result, err := r.coll.InsertOne(context.Background(), user)
	if err != nil {
		utils.Logger.WithField("error", err).Error("Failed to save user")

		return err
	}
	utils.Logger.WithField("user_id", result.InsertedID).Info("User added successfully")
	return nil

}

func (r Repo) Get(id string) (*models.User, error) {
	objectID, err := r.getObjectID(id)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("while get objectHex: %v", err))
		return nil, err
	}
	utils.Logger.WithField("action", "get_user").Info("Fetching user by ID")
	var user models.User
	err = r.coll.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		utils.Logger.WithField("id", id).Warn("User not found")
		return nil, err
	}
	return &user, nil
}
func (r Repo) GetAll() ([]models.User, error) {
	utils.Logger.WithField("action", "fetch_users").Info("Fetching all users")
	cursor, err := r.coll.Find(context.Background(), bson.D{})
	if err != nil {
		utils.Logger.WithField("error", err).Error("Failed to fetch users")
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []models.User
	if err = cursor.All(context.Background(), &users); err != nil {
		utils.Logger.WithField("error", err).Error("Failed to decode users")

		return nil, err
	}
	utils.Logger.WithField("user_count", len(users)).Info("Users fetched successfully")

	return users, nil
}

func (r Repo) Update(user *models.User) (err error) {
	objectID, err := r.getObjectID(user.ID.String())
	utils.Logger.WithField("action", "update_user").Info("Updating user by ID")

	opts := options.Update().SetUpsert(true)
	user.ID = objectID
	updateFilter := bson.M{"$set": bson.D{
		{"nickname", user.Nickname},
		{"password", user.Password},
	}}
	_, err = r.coll.UpdateOne(
		context.Background(), bson.M{"_id": objectID}, updateFilter,
		opts,
	)
	if err != nil {
		utils.Logger.WithField("error", err).Error("Error updating user")
		return err

	}
	utils.Logger.WithField("id", user.ID).Info("User updated successfully")
	return nil

}

func (r Repo) Delete(id string) error {
	objectID, err := r.getObjectID(id)
	utils.Logger.WithField("action", "delete_user").Info("Deleting user by ID")
	_, err = r.coll.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		utils.Logger.WithField("error", err).Error("Error deleting user")
		return nil
	}
	utils.Logger.WithField("id", id).Info("User deleted successfully")
	return nil
}
func (r *Repo) UserFilter(Nickname string) ([]models.User, error) {
	utils.Logger.WithField("action", "fetch_users_by_filter").Info("Fetching all users by filter")

	// Проверка на пустое значение или пробелы
	if strings.TrimSpace(Nickname) == "" {
		return nil, errors.New("nickname cannot be empty or whitespace")
	}

	// Создаем контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Фильтр с учётом регистра (регулярное выражение)
	filter := bson.M{"nickname": bson.M{"$regex": Nickname, "$options": "i"}}
	utils.Logger.WithField("filter", filter).Info("Applying filter to MongoDB query")

	// Выполняем запрос
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Обработка результата
	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Проверяем ошибки курсора
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	utils.Logger.WithField("count", len(users)).Info("Users fetched successfully")
	return users, nil
}
func (r *Repo) UserSorting(typeofsorting string) ([]models.User, error) {
	utils.Logger.WithField("action", "fetch_users_by_sorting").Info("Fetching all users by sorting")

	// Настройка сортировки
	var sortOrder int
	if typeofsorting == "asc" {
		sortOrder = 1 // По возрастанию
	} else if typeofsorting == "desc" {
		sortOrder = -1 // По убыванию
	} else {
		return nil, errors.New("invalid sorting type: must be 'asc' or 'desc'")
	}

	// Параметры поиска с сортировкой
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"nickname", sortOrder}})

	// Контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Выполнение запроса к MongoDB
	cursor, err := r.coll.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to fetch users")
		return nil, err
	}
	defer cursor.Close(ctx)

	// Обработка результатов
	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			utils.Logger.WithError(err).Error("Failed to decode user")
			return nil, err
		}
		users = append(users, user)
	}

	// Проверка на ошибки курсора
	if err := cursor.Err(); err != nil {
		utils.Logger.WithError(err).Error("Cursor error during iteration")
		return nil, err
	}

	utils.Logger.WithField("count", len(users)).Info("Users fetched successfully")
	return users, nil
}

func (r *Repo) UserPaging(skip int, limit int) ([]models.User, error) {
	utils.Logger.WithField("action", "fetch_users_by_paging").Info("Fetching all users by paging")
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.D{{"createdAt", 1}}) // Сортировка по дате создания

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.coll.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		utils.Logger.WithField("error", err).Error("Failed to fetch users")
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		utils.Logger.WithField("error", err).Error("Failed to decode users")
		return nil, err
	}

	return users, nil
}

func (r Repo) getObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
