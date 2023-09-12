package clusterclient

import (
	"bytes"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"

	S "github.com/off-chain-storage/go-off-chain-storage/service"
	U "github.com/off-chain-storage/go-off-chain-storage/utils"
)

const fileSize16MB uint64 = 16000000

type FILE struct {
	filedata    []byte
	filename    string
	clusterinfo string
	dbname      string
}

// EXPOSE FUNCTION #3
func GetMongoClient(key string) (client *mongo.Client, ctx context.Context, cancel context.CancelFunc, error error) {
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)

	// MongoDB_STORAGE
	uri := S.QuerySmartContract(key)
	opts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, opts)
	U.CheckErr(err)

	return client, ctx, cancel, err
}

func (f *FILE) HandleUploadRequest() {
	if f.clusterinfo == "Mongo Storage" {
		uploadforMongoStorage(f)
	} else {
		// uploadforIPFSCluster(f)
	}
}

func uploadforMongoStorage(f *FILE) {
	client, ctx, cancel, err := GetMongoClient(f.clusterinfo)
	U.CheckErr(err)

	defer client.Disconnect(ctx)
	defer cancel()

	db := client.Database(f.dbname)

	if uint64(len(f.filedata)) > fileSize16MB {
		bucket, err := gridfs.NewBucket(db)
		U.CheckErr(err)

		objectID, err := bucket.UploadFromStream(f.filename, f.filedata)
		U.CheckErr(err)
	} else {
		collection := db.Collection("fs.files")
		U.CheckErr(err)

		binData := primitive.Binary{
			Subtype: 0x00,
			Data:    f.filedata,
		}

		fileDoc := bson.D{
			{"filename", f.filename},
			{"length", uint64(len(f.filedata))},
			{"data", binData},
		}

		_, err = collection.InsertOne(context.Background(), fileDoc)
		U.CheckErr(err)
	}
}

func (f *FILE) HandleDownloadRequest() []byte {
	if f.clusterinfo == "Mongo Storage" {
		return downloadforMongoStorage(f)
	} else {
		// return downloadforIPFSCluster(f)

		// temp code
		var fileData []byte
		return fileData
	}
}

func downloadforMongoStorage(f *FILE) []byte {
	client, ctx, cancel, err := GetMongoClient(f.clusterinfo)
	U.CheckErr(err)

	defer client.Disconnect(ctx)
	defer cancel()

	db := client.Database(f.dbname)
	fsFiles := db.Collection("fs.files")

	filter := bson.D{{"filename", f.filename}}

	var result bson.M
	err = fsFiles.FindOne(ctx, filter).Decode(&result)
	U.CheckErr(err)

	if result["length"].(int64) > fileSize16MB {
		return downloadLargeFile(f.filename, db)
	} else {
		return downloadSmallFile(result)
	}
}

func downloadSmallFile(result bson.M) []byte {
	dataBinData, _ := result["data"].(primitive.Binary)
	data := dataBinData.Data
	return data
}

func downloadLargeFile(filename string, db *mongo.Database) []byte {
	bucket, _ := gridfs.NewBucket(db)

	var buf bytes.Buffer
	dStream, err := bucket.DownloadToStreamByName(filename, &buf)
	U.CheckErr(err)

	return buf.Bytes()
}
