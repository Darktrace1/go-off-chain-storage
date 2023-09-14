package clusterclient

import (
	"bytes"
	"context"
	"fmt"
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
	Filedata    []byte
	Filename    string
	Clusterinfo string
	Dbname      string
	Dbdata      bson.M
}

// EXPOSE FUNCTION #3
func GetMongoClient() (client *mongo.Client, ctx context.Context, cancel context.CancelFunc, error error) {
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)

	// MongoDB_STORAGE
	uri := S.QuerySmartContract("MongoDB")

	opts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, opts)
	U.CheckErr(err)

	return client, ctx, cancel, err
}

func (f *FILE) HandleUploadRequest() {
	if f.Clusterinfo == "MongoDB Storage" {
		uploadforMongoStorage(f)
	} else if f.Clusterinfo == "MongoDB" {
		// usingMongoDB(f)
		fmt.Println("HI")
	} else {
		// uploadforIPFSCluster(f)
		fmt.Println("HI")
	}
}

func uploadforMongoStorage(f *FILE) {
	client, ctx, cancel, err := GetMongoClient()
	U.CheckErr(err)

	defer client.Disconnect(ctx)
	defer cancel()

	db := client.Database(f.Dbname)

	if uint64(len(f.Filedata)) > fileSize16MB {
		bucket, err := gridfs.NewBucket(db)
		U.CheckErr(err)

		_, err = bucket.UploadFromStream(f.Filename, bytes.NewReader(f.Filedata))
		U.CheckErr(err)
	} else {
		collection := db.Collection("fs.files")
		U.CheckErr(err)

		fileDoc := bson.M{
			"filename": f.Filename,
			"length":   uint64(len(f.Filedata)),
			"data": primitive.Binary{
				Subtype: 0x00,
				Data:    f.Filedata,
			},
		}

		_, err = collection.InsertOne(context.Background(), fileDoc)
		U.CheckErr(err)
	}
}

func (f *FILE) HandleDownloadRequest() []byte {
	if f.Clusterinfo == "Mongo Storage" {
		return downloadforMongoStorage(f)
	} else {
		// return downloadforIPFSCluster(f)

		// temp code
		var FileData []byte
		return FileData
	}
}

func downloadforMongoStorage(f *FILE) []byte {
	client, ctx, cancel, err := GetMongoClient()
	U.CheckErr(err)

	defer client.Disconnect(ctx)
	defer cancel()

	db := client.Database(f.Dbname)
	fsFiles := db.Collection("fs.files")

	filter := bson.M{"filename": f.Filename}

	var result bson.M
	err = fsFiles.FindOne(ctx, filter).Decode(&result)
	U.CheckErr(err)

	if result["length"].(uint64) > fileSize16MB {
		return downloadLargeFile(f.Filename, db)
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
	_, err := bucket.DownloadToStreamByName(filename, &buf)
	U.CheckErr(err)

	return buf.Bytes()
}
