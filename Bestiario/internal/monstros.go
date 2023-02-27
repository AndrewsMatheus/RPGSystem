package monstros

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func init() {

}

type Monstro struct {
	ID        uint
	Nome      string
	Descricao string
	UUID      uuid.UUID
}

type MonstrosRepository interface {
	GetMonstroById(ctx context.Context, ID uuid.UUID) (*Monstro, error)
	SearchMonstros(ctx context.Context, nome string) (Monstro, int, error)
	SaveMonstro(ctx context.Context, m Monstro) (*Monstro, error)
	UpdateMonstro(ctx context.Context, m *Monstro)
	DeleteMonstro(ctx context.Context, ID uuid.UUID) (*Monstro, error)
	Open() error
	Close() error
}

type MonstroGorm struct {
	ID        uint   `gorm:"primaryKey;column:ID"`
	UUID      string `gorm:"uniqueIndex;column:UUID"`
	Nome      string `gorm:"column:NOME"`
	Descricao string `gorm:"column:"`
}

func (m MonstroGorm) ToEntity() (Monstro, error) {
	parsed, err := uuid.Parse(m.UUID)
	if err != nil {
		return Monstro{}, err
	}

	return Monstro{
		ID:        m.ID,
		UUID:      parsed,
		Nome:      m.Nome,
		Descricao: m.Descricao,
	}, nil
}

type MonstroRepository struct {
	connection *gorm.DB
	dsn        string
}

func NewMonstroRepository(dsn string) *MonstroRepository {
	var newRepository MonstroRepository
	newRepository.dsn = dsn

	return &newRepository
}

func (repository *MonstroRepository) Open() error {
	var err error
	repository.connection, err = gorm.Open(mysql.Open(repository.dsn), &gorm.Config{})

	if err != nil {
		return err
	}
	return nil
}
func (repository *MonstroRepository) GetMonstro(ctx context.Context, ID uuid.UUID) (*Monstro, error) {
	var row MonstroGorm
	err := repository.connection.WithContext(ctx).Table("monstros").First(&row, "uuid = ?", ID).Error
	if err != nil {
		return nil, err
	}

	monstro, err := row.ToEntity()
	if err != nil {
		return nil, err
	}

	return &monstro, nil
}
func (repository *MonstroRepository) SaveMonstro(ctx context.Context, monstro Monstro) (*Monstro, error) {
	row := NewRow(monstro)
	err := repository.connection.WithContext(ctx).Table("monstros").Save(&row).Error
	if err != nil {
		return nil, err
	}

	monstro, err = row.ToEntity()
	if err != nil {
		return nil, err
	}

	return &monstro, nil
}

func NewRow(monstro Monstro) MonstroGorm {
	return MonstroGorm{
		UUID:      uuid.NewString(),
		Nome:      monstro.Nome,
		Descricao: monstro.Descricao,
	}
}

func (repository *MonstroRepository) CreateMonstro(ctx context.Context, m Monstro) (*Monstro, error) {
	tx := repository.connection.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	var total int64
	var err error

	if err != nil {
		tx.Rollback()
		return nil, err
	} else if total > 0 {
		tx.Rollback()
		return nil, errors.New("já existe um registro deste no banco")
	}

	var row = NewRow(m)

	err = tx.Save(&row).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &m, nil
}

func (repository *MonstroRepository) SearchMonstros(ctx context.Context, nome string) (*Monstro, error) {
	var row MonstroGorm
	err := repository.connection.WithContext(ctx).Table("monstros").Where("NOME LIKE ?", "%"+nome+"%").Find(&row).Error
	if err != nil {
		return nil, err

	}

	monstro, err := row.ToEntity()
	if err != nil {
		return nil, err
	}

	return &monstro, nil
}

func (repository *MonstroRepository) UpdateMonstro(ctx context.Context, m *Monstro) []Monstro {
	var monstroAtualizado []Monstro
	repository.connection.Model(&monstroAtualizado).Clauses(clause.Returning{}).Where("ID = ?", m.ID).Updates(m)
	return monstroAtualizado
}
func (repository *MonstroRepository) DeleteMonstro(ctx context.Context, m *Monstro) []Monstro {
	var monstroDeletado []Monstro

	repository.connection.Clauses(clause.Returning{}).Where("ID = ?", m.ID).Delete(&monstroDeletado)
	return monstroDeletado
}
