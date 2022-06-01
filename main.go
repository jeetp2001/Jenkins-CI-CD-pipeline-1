package main



import (

	"net/http"

	"html/template"

	"os"

	"github.com/aws/aws-sdk-go/aws"

    "github.com/aws/aws-sdk-go/aws/session"

    "github.com/aws/aws-sdk-go/service/ec2"

	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"log"

	// "reflect"

	"strconv"

)



var tpl = template.Must(template.ParseFiles("templates/index.html"))

var tpl1 = template.Must(template.ParseFiles("templates/ec2fun.html"))

var tpl2 = template.Must(template.ParseFiles("templates/ec2done.html"))

var tpl3 = template.Must(template.ParseFiles("templates/ec2showins.html"))

var tpl4 = template.Must(template.ParseFiles("templates/stopec2ins.html"))

var tpl5 = template.Must(template.ParseFiles("templates/delec2ins.html"))

var tpl6 = template.Must(template.ParseFiles("templates/s3fun.html"))

var tpl7 = template.Must(template.ParseFiles("templates/s3bucket.html"))

var tpl8 = template.Must(template.ParseFiles("templates/s3error.html"))

var tpl9 = template.Must(template.ParseFiles("templates/lists3.html"))

var tpl10 = template.Must(template.ParseFiles("templates/delbuc.html"))

var tpl11 = template.Must(template.ParseFiles("templates/delbucerror.html"))

var tpl12 = template.Must(template.ParseFiles("templates/delallitems.html"))



func indexHandler(w http.ResponseWriter, r *http.Request){

	tpl.Execute(w, nil);

}



func ec2fun(w http.ResponseWriter, r *http.Request){

	tpl1.Execute(w, nil);

}



func ec2t2micro(w http.ResponseWriter, r *http.Request){

	noi, err := strconv.ParseInt(r.FormValue("noi"),10,64)

    if err != nil {

        fmt.Println(err)

    }

	sess, err := session.NewSession(&aws.Config{

		Region: aws.String("ap-south-1"),

	})

	svc := ec2.New(sess)

	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{



		ImageId: aws.String("ami-079b5e5b3971bd10d"),

		InstanceType: aws.String("t2.micro"),

		MinCount: aws.Int64(noi),

		MaxCount: aws.Int64(noi),

	})



	if err != nil {

		fmt.Println("Could not create instance ", err)

	}



	fmt.Println("created instance ", *runResult.Instances[0].InstanceId)



	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{

		Resources: []*string{runResult.Instances[0].InstanceId},

		Tags: []*ec2.Tag{

			{

				Key: aws.String("Name"),

				Value: aws.String("ins-by-go"),

			},

		},

	})

	if errtag != nil {

		log.Println("Could not create tag ", runResult.Instances[0].InstanceId, errtag)

	}



	fmt.Println("successfully tagged instance")

	tpl2.Execute(w, nil);

}



func ec2showins(w http.ResponseWriter, r *http.Request){

	sess, err := session.NewSession(&aws.Config{

		Region: aws.String("ap-south-1"),

	})

	ec2svc := ec2.New(sess)

	result, err := ec2svc.DescribeInstances(nil)

	if err != nil {

		fmt.Println("Error", err)

	} else {

		// for idx,res := range result.Reservations {

		// 	for _,inst := range result.Reservations[idx].Instances {

		// 		// fmt.Println("   - ", *inst.InstanceId)

		// 		// fmt.Println("   - ", *inst.State.Name)

		// 		// fmt.Println(" ................... ")

		// 		// insname := append(insname, *inst.InstanceId)

		// 		// fmt.Println(insname)

		// 		// insname := append(insname, inst)

		// 	}

		// }

		tpl3.Execute(w, result.Reservations)

	}

	fmt.Println(result.Reservations)

}



func printfun(w http.ResponseWriter, r *http.Request){

	fmt.Println("done!!!!!")

}



func stopec2ins(w http.ResponseWriter, r *http.Request){

	sess, err := session.NewSession(&aws.Config{

		Region: aws.String("ap-south-1"),

	})

	if err != nil {

		fmt.Println(err)

	}

	ec2svc := ec2.New(sess)



	input := &ec2.StopInstancesInput{

		InstanceIds: []*string{

			aws.String(r.FormValue("insid")),

		},

		DryRun: aws.Bool(false),

	}

	runResult, err := ec2svc.StopInstances(input)

	fmt.Println(w, runResult)

	tpl4.Execute(w, nil);

}



func delec2ins(w http.ResponseWriter, r *http.Request){

	sess, err := session.NewSession(&aws.Config{

		Region: aws.String("ap-south-1"),

	})

	if err != nil {

		fmt.Println(err)

	}

	ec2svc := ec2.New(sess)



	input := &ec2.TerminateInstancesInput{

		InstanceIds: []*string{

			aws.String(r.FormValue("insid")),

		},

		DryRun: aws.Bool(false),

	}

	runResult, err := ec2svc.TerminateInstances(input)

	fmt.Println(w, runResult)

	tpl5.Execute(w, nil);

}



func s3fun(w http.ResponseWriter, r *http.Request){

	tpl6.Execute(w, nil)

}



func s3bucket(w http.ResponseWriter, r *http.Request){

	sess, err := session.NewSession(&aws.Config{

        Region: aws.String("us-west-2")},

    )

	svc := s3.New(sess)

	bucname := r.FormValue("bucname")

	input := &s3.CreateBucketInput{

		Bucket: aws.String(bucname),

	}

	result, err := svc.CreateBucket(input)

	if err != nil {

		if aerr, ok := err.(awserr.Error); ok {

			switch aerr.Code() {

			case s3.ErrCodeBucketAlreadyExists:

				fmt.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())

				tpl8.Execute(w, nil)

			case s3.ErrCodeBucketAlreadyOwnedByYou:

				fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())

				tpl8.Execute(w, nil)

			default:

				fmt.Println(aerr.Error())

			}

		} else {

			fmt.Println(err.Error())

		}

		return

	}		

	

	fmt.Println(result)

	tpl7.Execute(w, nil)

}



func lists3(w http.ResponseWriter, r *http.Request){

	sess, err := session.NewSession(&aws.Config{

		Region: aws.String("ap-south-1"),

	})

	s3svc := s3.New(sess)

	

	input := &s3.ListBucketsInput{}



	result, err := s3svc.ListBuckets(input)

	if err != nil {

		if aerr, ok := err.(awserr.Error); ok {

			switch aerr.Code() {

			default:

				fmt.Println(aerr.Error())

			}

		} else {

			fmt.Println(err.Error())

		}

		return

	}

	tpl9.Execute(w, result.Buckets)

}

func exitErrorf(msg string, args ...interface{}) {

	fmt.Fprintf(os.Stderr, msg+"\n", args...)

	os.Exit(1)

}



func delallitems(w http.ResponseWriter, r *http.Request){

	sess, _ := session.NewSession(&aws.Config{

		Region: aws.String("ap-south-1")},

	)

	svc := s3.New(sess)

	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{

		Bucket: aws.String(r.FormValue("s3name")),

	})

	

	if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {

		fmt.Println("Unable to delete objects from bucket %q, %v", r.FormValue("s3name"), err)

	}

	tpl12.Execute(w, nil)

}





func delbuc(w http.ResponseWriter, r *http.Request){

	sess, err := session.NewSession(&aws.Config{

		Region: aws.String("ap-south-1"),

	})

	if err != nil {

		fmt.Println(err)

	}

	

	s3svc := s3.New(sess)

	_, err = s3svc.DeleteBucket(&s3.DeleteBucketInput{

		Bucket: aws.String(r.FormValue("s3name")),

	})

	if err != nil {

		fmt.Println("Unable to delete bucket %q, %v", r.FormValue("s3name"), err)

		tpl11.Execute(w, err)

	}

	tpl10.Execute(w, nil)



	// err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{

	// 	Bucket: aws.String(bucket),

	// })

}



func main(){

	mux := http.NewServeMux()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")

	if port == "" {

		port = "3000"

	}



	mux.HandleFunc("/", indexHandler)

	mux.HandleFunc("/ec2fun", ec2fun)

	mux.HandleFunc("/ec2t2micro", ec2t2micro)

	mux.HandleFunc("/ec2showins", ec2showins)

	mux.HandleFunc("/stopec2ins", stopec2ins)

	mux.HandleFunc("/delec2ins", delec2ins)

	mux.HandleFunc("/s3fun", s3fun)

	mux.HandleFunc("/s3bucket", s3bucket)

	mux.HandleFunc("/lists3", lists3)

	mux.HandleFunc("/delbuc",delbuc)

	mux.HandleFunc("/delallitems",delallitems)

	http.ListenAndServe(":"+port, mux)

}
