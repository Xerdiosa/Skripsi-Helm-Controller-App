package repositories

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models"
	"gorm.io/gorm"
)

type KinesisProvider struct {
	database *gorm.DB
	kinesis  *kinesis.Client
}

func InitKinesisProvider(db *gorm.DB, kinesis *kinesis.Client) Providers {
	kinesisProvider := &KinesisProvider{}
	kinesisProvider.database = db
	kinesisProvider.kinesis = kinesis

	return kinesisProvider
}

func (k *KinesisProvider) Convert(rawData interface{}) (interface{}, error) {
	jsonStr, err := json.Marshal(rawData)
	if err != nil {
		return nil, err
	}
	component := models.Kinesis{}
	err = json.Unmarshal(jsonStr, &component)
	if err != nil {
		return nil, err
	}
	return component, nil
}

func (k *KinesisProvider) PreProcess(data interface{}, prevData interface{}, module interface{}, moduleRelease interface{}) (interface{}, error) {
	processed, ok := data.(models.Kinesis)
	if !ok {
		err := errors.New("conversion to kinesis failed")
		return nil, err
	}
	releaseParsed, ok := moduleRelease.(models.ModuleRelease)
	if !ok {
		err := errors.New("conversion to release failed")
		return nil, err
	}

	var oldData models.Kinesis
	if prevData != nil {
		oldData, ok = prevData.(models.Kinesis)
		if !ok {
			err := errors.New("conversion to kinesis failed")
			return nil, err
		}
	}

	processed.ModuleReleaseID = releaseParsed.ID
	processed.Revision = oldData.Revision + 1
	return processed, nil
}

func (k *KinesisProvider) InstallComponent(kinesisInterface interface{}) error {
	kinesisData, ok := kinesisInterface.(models.Kinesis)
	if !ok {
		err := errors.New("conversion to kinesis failed")
		return err
	}

	input := kinesis.CreateStreamInput{
		StreamName: &kinesisData.Name,
		ShardCount: &kinesisData.Shards,
	}
	_, err := k.kinesis.CreateStream(context.TODO(), &input)
	return err
}

func (k *KinesisProvider) UpdateComponent(kinesisInterface interface{}) error {
	kinesisData, ok := kinesisInterface.(models.Kinesis)
	if !ok {
		err := errors.New("conversion to kinesis failed")
		return err
	}

	input := kinesis.UpdateShardCountInput{
		ScalingType:      types.ScalingTypeUniformScaling,
		StreamName:       &kinesisData.Name,
		TargetShardCount: &kinesisData.Shards,
	}

	_, err := k.kinesis.UpdateShardCount(context.TODO(), &input)
	return err
}

func (k *KinesisProvider) UninstallComponent(kinesisInterface interface{}) error {
	kinesisData, ok := kinesisInterface.(models.Kinesis)
	if !ok {
		err := errors.New("conversion to kinesis failed")
		return err
	}
	input := kinesis.DeleteStreamInput{
		StreamName: &kinesisData.Name,
	}
	_, err := k.kinesis.DeleteStream(context.TODO(), &input)
	return err

}

func (k *KinesisProvider) GetAllName() ([]string, error) {
	var names []string
	result := k.database.Model(&models.Kinesis{}).Pluck("name", &names)
	if result.Error != nil {
		return nil, result.Error
	}
	return names, nil
}
func (k *KinesisProvider) Add(kinesisInterface interface{}) error {
	kinesis, ok := kinesisInterface.(models.Kinesis)
	if !ok {
		err := errors.New("conversion to kinesis failed")
		return err
	}

	result := k.database.Create(&kinesis)
	return result.Error
}
func (k *KinesisProvider) Remove(kinesisInterface interface{}) error {
	kinesis, ok := kinesisInterface.(models.Kinesis)
	if !ok {
		err := errors.New("conversion to kinesis failed")
		return err
	}

	result := k.database.Delete(&models.Kinesis{}, "name = ?", kinesis.Name)
	return result.Error
}

func (k *KinesisProvider) Update(kinesisInterface interface{}) error {
	kinesis, ok := kinesisInterface.(models.Kinesis)
	if !ok {
		err := errors.New("conversion to kinesis failed")
		return err
	}

	result := k.database.Model(&kinesis).Where("name = ?", kinesis.Name).Updates(kinesis)
	return result.Error
}
func (k *KinesisProvider) GetDetail(releaseName string) (interface{}, error) {
	var kinesis models.Kinesis
	result := k.database.Where("name = ?", releaseName).First(&kinesis)
	return kinesis, result.Error
}

func (k *KinesisProvider) GetDetailFromComponent(kinesisInterface interface{}) (interface{}, error) {
	kinesis, ok := kinesisInterface.(models.Kinesis)
	if !ok {
		err := errors.New("conversion to Kinesis failed")
		return nil, err
	}

	return k.GetDetail(kinesis.Name)
}

func (k *KinesisProvider) GetFromModuleReleaseID(ModuleReleaseID uint) ([]interface{}, error) {
	var kinesis []models.Kinesis
	result := k.database.Where("module_release_id = ?", ModuleReleaseID).Find(&kinesis)

	var kinesisInterface []interface{} = make([]interface{}, len(kinesis))
	for i, v := range kinesis {
		kinesisInterface[i] = v
	}

	return kinesisInterface, result.Error
}
