package services

import (
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/repositories"
	"gorm.io/gorm"
)

type IKinesisService interface {
	InstallOrUpgradeKinesis(models.Kinesis) error
	GetAllReleaseName() ([]string, error)
	GetReleaseDetail(string) (models.Kinesis, error)
	RemoveKinesis(string) error
}

type KinesisService struct {
	kinesisProvider repositories.Providers
}

func InitKinesisService(kinesisProvider repositories.Providers) IKinesisService {
	KinesisService := &KinesisService{}
	KinesisService.kinesisProvider = kinesisProvider
	return KinesisService
}

func (k *KinesisService) InstallOrUpgradeKinesis(kinesis models.Kinesis) error {
	oldKinesisInterface, err := k.kinesisProvider.GetDetail(kinesis.Name)
	oldKinesis := oldKinesisInterface.(models.Kinesis)
	if err == gorm.ErrRecordNotFound {
		return k.installKinesis(kinesis)
	}
	return k.upgradeKinesis(kinesis, oldKinesis)
}

func (k *KinesisService) installKinesis(kinesis models.Kinesis) error {
	kinesis.Revision = 1
	err := k.kinesisProvider.InstallComponent(kinesis)
	if err != nil {
		return err
	}
	err = k.kinesisProvider.Add(kinesis)
	return err
}

func (k *KinesisService) upgradeKinesis(kinesis models.Kinesis, oldKinesis models.Kinesis) error {
	kinesis.Revision = oldKinesis.Revision + 1

	err := k.kinesisProvider.UpdateComponent(kinesis)
	if err != nil {
		return err
	}

	err = k.kinesisProvider.Update(kinesis)
	return err
}

func (k *KinesisService) RemoveKinesis(kinesis string) error {
	kinesisInstance, err := k.kinesisProvider.GetDetail(kinesis)
	if err != nil {
		return err
	}
	err = k.kinesisProvider.UninstallComponent(kinesisInstance)
	if err != nil {
		return err
	}
	err = k.kinesisProvider.Remove(kinesisInstance)
	return err
}

func (k *KinesisService) GetAllReleaseName() ([]string, error) {
	result, err := k.kinesisProvider.GetAllName()
	return result, err
}

func (k *KinesisService) GetReleaseDetail(releaseName string) (models.Kinesis, error) {
	resultInterface, err := k.kinesisProvider.GetDetail(releaseName)
	result := resultInterface.(models.Kinesis)
	return result, err
}
