package voiceit2

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	voiceit2 "github.com/voiceittech/voiceit2-go"
)

// Test API Key and Token in environment variables
func TestKeyToken(t *testing.T) {
	assert := assert.New(t)
	apikey := os.Getenv("VIAPIKEY")
	apitoken := os.Getenv("VIAPITOKEN")
	assert.NotEqual("", apikey, "Please set and export VIAPIKEY environment variable in order to test")
	assert.NotEqual("", apitoken, "Please set and export VIAPITOKEN environment variable in order to test")
}

// GetAllUsers()
type GetAllUsersReturn struct {
	Message      string       `json:"message"`
	Count        int          `json:"count"`
	Status       int          `json:"status"`
	TimeTaken    string       `json:"timeTaken"`
	Users        []GetAllUser `json:"users"`
	ResponseCode string       `json:"responseCode"`
}

type GetAllUser struct {
	CreatedAt int    `json:"createdAt"`
	UserId    string `json:"userId"`
}

// Try getting all users
func TestGetAllUsers(t *testing.T) {
	assert := assert.New(t)
	apikey := os.Getenv("VIAPIKEY")
	apitoken := os.Getenv("VIAPITOKEN")
	myVoiceIt := voiceit2.NewClient(apikey, apitoken)
	var gaur GetAllUsersReturn
	err := json.Unmarshal([]byte(myVoiceIt.GetAllUsers()), &gaur)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NotEqual("", gaur.Message, "message return from GetAllUsers() call is empty")
	assert.Equal(200, gaur.Status, "status return from GetAllUsers() call is not 200")
	assert.NotEqual("", gaur.TimeTaken, "timeTaken return from GetAllUsers() call is empty")
	assert.Equal(gaur.Count, len(gaur.Users), "count return from GetAllUsers() does not match the length of the returned number of users")
	for _, elem := range gaur.Users {
		assert.NotEqual(0, elem.CreatedAt, "createdAt for user did not return a real date/time integer")
		assert.NotNil(elem.UserId, "userId return from GetAllUsers() call is empty")
	}
	assert.Equal("SUCC", gaur.ResponseCode, "responseCode return from GetAllUsers() call not success")
}

// CreateUser() /DeleteUser

type CreateUserReturn struct {
	Message      string `json:"message"`
	Status       int    `json:"status"`
	TimeTaken    string `json:"timeTaken"`
	CreatedAt    int    `json:"createdAt"`
	UserId       string `json:"userId"`
	ResponseCode string `json:"responseCode"`
}

type DeleteUserReturn struct {
	Message      string `json:"message"`
	Status       int    `json:"status"`
	TimeTaken    string `json:"timeTaken"`
	ResponseCode string `json:"responseCode"`
}

// Try creating user and deleting same user
func TestCreateUserDeleteUser(t *testing.T) {
	assert := assert.New(t)
	apikey := os.Getenv("VIAPIKEY")
	apitoken := os.Getenv("VIAPITOKEN")
	myVoiceIt := voiceit2.NewClient(apikey, apitoken)

	// CreateUser()
	var cur CreateUserReturn
	err1 := json.Unmarshal([]byte(myVoiceIt.CreateUser()), &cur)
	if err1 != nil {
		t.Error(err1.Error())
	}

	assert.NotEqual("", cur.Message, "message return from CreateUser() call is empty")
	assert.Equal("Created user with userId : usr_", cur.Message[0:31], "message return from CreateUser() does not follow the pattern \"Created user with userId: usr_00000000000000000000000000000000\"")
	assert.Equal(201, cur.Status, "status return from CreateUser() call is not 201")
	assert.NotEqual("", cur.TimeTaken, "timeTaken return from CreateUser() call is empty")
	assert.NotEqual(0, cur.CreatedAt, "createdAt for user did not return a real date/time integer")
	assert.NotNil(cur.UserId, "userId return from CreateUser() call is empty")
	assert.Equal(36, len(cur.UserId), "userId return from CreateUser() not a string of length 36")
	assert.Equal("usr_", string(cur.UserId[0:4]), "userId return from CreateUser() does not follow the convention \"usr_00000000000000000000000000000000\"")
	assert.Equal("SUCC", cur.ResponseCode, "responseCode return from CreateUser() is not \"SUCC\"")

	// DeleteUser()
	var dur DeleteUserReturn
	err2 := json.Unmarshal([]byte(myVoiceIt.DeleteUser(cur.UserId)), &dur)
	if err2 != nil {
		t.Error(err2.Error())
	}
	assert.NotEqual("", dur.Message, "message return from DeleteUser() call is empty")
	assert.Equal("Deleted user with userId : usr_", dur.Message[0:31], "message return from DeleteUser() does not follow the convention \"Deleted user with userId : usr_00000000000000000000000000000000\"")
	assert.Equal(200, dur.Status, "status return from CreateUser() call is not 201")
	assert.Equal("SUCC", dur.ResponseCode, "responseCode return from DeleteUser() is not \"SUCC\"")
}

type CreateGroupReturn struct {
	Message      string `json:"message"`
	Description  string `json:"description"`
	GroupId      string `json:"groupId"`
	Status       int    `json:"status"`
	CreatedAt    int    `json:"createdAt"`
	TimeTaken    string `json:"timeTaken"`
	ResponseCode string `json:"responseCode"`
}

type AddUserToGroupReturn struct {
	Message      string `json:"message"`
	Status       int    `json:"status"`
	TimeTaken    string `json:"timeTaken"`
	ResponseCode string `json:"responseCode"`
}

type RemoveUserFromGroupReturn struct {
	Message      string `json:"message"`
	Status       int    `json:"status"`
	TimeTaken    string `json:"timeTaken"`
	ResponseCode string `json:"responseCode"`
}

type DeleteGroupReturn struct {
	Message      string `json:"message"`
	Status       int    `json:"status"`
	TimeTaken    string `json:"timeTaken"`
	ResponseCode string `json:"responseCode"`
}

// Test create Group, Create User, Add user to that group, remove user from group, delete user, delete group
func TestCreateUserGroupInteractions(t *testing.T) {
	assert := assert.New(t)
	apikey := os.Getenv("VIAPIKEY")
	apitoken := os.Getenv("VIAPITOKEN")
	myVoiceIt := voiceit2.NewClient(apikey, apitoken)

	// CreateUser()
	var cur CreateUserReturn
	err1 := json.Unmarshal([]byte(myVoiceIt.CreateUser()), &cur)
	if err1 != nil {
		t.Error(err1.Error())
	}
	userId := cur.UserId

	// CreateGroup()
	var cgr CreateGroupReturn
	err2 := json.Unmarshal([]byte(myVoiceIt.CreateGroup("Sample Group Description")), &cgr)
	if err2 != nil {
		t.Error(err2.Error())
	}
	groupId := cgr.GroupId
	assert.Equal(64, len(cgr.Message), "message return from CreateGroup() does not follow the pattern \"Created group with groupId: grp_00000000000000000000000000000000\"")
	assert.Equal("Created group with groupId: grp_", string(cgr.Message[0:32]), "message return from CreateGroup() does not follow the pattern \"Created group with groupId: grp_00000000000000000000000000000000\"")
	assert.Equal(201, cgr.Status, "status return from CreateGroup() is not 201")
	assert.Equal("Sample Group Description", cgr.Description, "description return from CreateGroup() does not match passed value")
	assert.Equal("SUCC", cgr.ResponseCode, "responseCode return from CreateGroup() not \"SUCC\"")

	// AddUserToGroup()
	var autgr AddUserToGroupReturn
	err3 := json.Unmarshal([]byte(myVoiceIt.AddUserToGroup(groupId, userId)), &autgr)
	if err3 != nil {
		t.Error(err3.Error())
	}
	r1, _ := regexp.Compile("Successfully added user usr_([a-z0-9]+){32} to group with groupId : grp_([a-z0-9]){32}")
	assert.True(r1.MatchString(autgr.Message), "message return from AddUserToGroup() does not follow the pattern \"Successfully added user usr_00000000000000000000000000000000 to group with groupId : grp_00000000000000000000000000000000\"")
	assert.Equal(200, autgr.Status, "status return from AddUserToGroup() is not 200")
	assert.Equal("SUCC", autgr.ResponseCode, "responseCode return from AddUserToGroup() not \"SUCC\"")

	// RemoveUserFromGroup()
	var rufgr RemoveUserFromGroupReturn
	err4 := json.Unmarshal([]byte(myVoiceIt.RemoveUserFromGroup(groupId, userId)), &rufgr)
	if err4 != nil {
		t.Error(err4.Error())
	}
	r2, _ := regexp.Compile("Successfully removed user usr_([a-z0-9]+){32} from group with groupId : grp_([a-z0-9]){32}")
	assert.True(r2.MatchString(rufgr.Message), "message return from RemoveUserFromGroup() does not follow the pattern \"Successfully removed user usr_00000000000000000000000000000000 to group with groupId : grp_00000000000000000000000000000000\"")
	assert.Equal(200, rufgr.Status, "status return from RemoveUserFromGroup() is not 200")
	assert.Equal("SUCC", rufgr.ResponseCode, "responseCode return from RemoveUserFromGroup() not \"SUCC\"")

	// DeleteGroup()
	var dgr DeleteGroupReturn
	err5 := json.Unmarshal([]byte(myVoiceIt.DeleteGroup(groupId)), &dgr)
	if err5 != nil {
		t.Error(err5.Error())
	}
	r3, _ := regexp.Compile("Successfully deleted group with groupId : grp_([a-z0-9]){32}")
	assert.True(r3.MatchString(dgr.Message), "message return from DeleteGroup() does not follow the pattern \"Successfully deleted group with groupId : grp_00000000000000000000000000000000\"")
	assert.Equal(200, dgr.Status, "status return from DeleteGroup() is not 200")
	assert.Equal("SUCC", dgr.ResponseCode, "responseCode return from DeleteGroup() not \"SUCC\"")

	myVoiceIt.DeleteUser(userId)
}

// Helper function to download files to disk
func downloadFromUrl(url string) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	log.Println("Downloading " + url + "...")
	output, err := os.Create(fileName)
	if err != nil {
		log.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		log.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		log.Println("Error while downloading", url, "-", err)
		return
	}

}

type CreateVideoEnrollmentReturn struct {
	Message         string  `json:"message"`
	ContentLanguage string  `json:"contentLanguage"`
	Id              int     `json:"id"`
	Status          int     `json:"status"`
	Text            string  `json:"text"`
	TextConfidence  float32 `json:"textConfidence"`
	CreatedAt       int     `json:"createdAt"`
	TimeTaken       string  `json:"timeTaken"`
	ResponseCode    string  `json:"responseCode"`
}

type VideoVerificationReturn struct {
	Message         string  `json:"message"`
	Status          int     `json:"status"`
	VoiceConfidence float32 `json:"voiceConfidence"`
	FaceConfidence  float32 `json:"faceConfidence"`
	Text            string  `json:"text"`
	TextConfidence  float32 `json:"textConfidence"`
	BlinksCount     int     `json:"blinksCount"`
	TimeTaken       string  `json:"timeTaken"`
	ResponseCode    string  `json:"responseCode"`
}

type VideoIdentificationReturn struct {
	Message         string  `json:"message"`
	UserId          string  `json:"userId"`
	Status          int     `json:"status"`
	VoiceConfidence float32 `json:"voiceConfidence"`
	FaceConfidence  float32 `json:"faceConfidence"`
	Text            string  `json:"text"`
	TextConfidence  float32 `json:"textConfidence"`
	BlinksCount     int     `json:"blinksCount"`
	TimeTaken       string  `json:"timeTaken"`
	ResponseCode    string  `json:"responseCode"`
}

type DeleteEnrollmentReturn struct {
	Message      string `json:"message"`
	Status       int    `json:"status"`
	TimeTaken    string `json:"timeTaken"`
	ResponseCode string `json:"responseCode"`
}

// Test video enrollment/verification/identification (and deleting each individually)
func TestVideoEnrollmentVerificationIdentification(t *testing.T) {
	assert := assert.New(t)
	apikey := os.Getenv("VIAPIKEY")
	apitoken := os.Getenv("VIAPITOKEN")
	myVoiceIt := voiceit2.NewClient(apikey, apitoken)

	// CreateUser()
	var cur CreateUserReturn
	err1 := json.Unmarshal([]byte(myVoiceIt.CreateUser()), &cur)
	if err1 != nil {
		t.Error(err1.Error())
	}
	userId := cur.UserId

	// Make 3 enrollments
	// Download 3 enrollment video
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoEnrollmentArmaan1.mov")
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoEnrollmentArmaan2.mov")
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoEnrollmentArmaan3.mov")
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoVerificationArmaan1.mov")

	// Enrollment1
	var cver1 CreateVideoEnrollmentReturn
	err2 := json.Unmarshal([]byte(myVoiceIt.CreateVideoEnrollment(userId, "en-US", "./videoEnrollmentArmaan1.mov", false)), &cver1)
	if err2 != nil {
		t.Error(err2.Error())
	}

	// Enrollment2
	var cver2 CreateVideoEnrollmentReturn
	err3 := json.Unmarshal([]byte(myVoiceIt.CreateVideoEnrollment(userId, "en-US", "./videoEnrollmentArmaan2.mov", false)), &cver2)
	if err3 != nil {
		t.Error(err3.Error())
	}

	// Enrollment3
	var cver3 CreateVideoEnrollmentReturn
	err4 := json.Unmarshal([]byte(myVoiceIt.CreateVideoEnrollment(userId, "en-US", "./videoEnrollmentArmaan3.mov", false)), &cver3)
	if err4 != nil {
		t.Error(err4.Error())
	}

	// Run checks on enrollment returns
	r1, _ := regexp.Compile("Successfully added video enrollment for user with userId : usr_([a-z0-9]){32}")
	assert.True(r1.MatchString(cver1.Message), "message return from CreateVideoEnrollment() does not follow the pattern \"Successfully added video enrollment for user with userId : usr_00000000000000000000000000000000\"")
	assert.True(r1.MatchString(cver2.Message), "message return from CreateVideoEnrollment() does not follow the pattern \"Successfully added video enrollment for user with userId : usr_00000000000000000000000000000000\"")
	assert.True(r1.MatchString(cver3.Message), "message return from CreateVideoEnrollment() does not follow the pattern \"Successfully added video enrollment for user with userId : usr_00000000000000000000000000000000\"")
	assert.Equal("en-US", cver1.ContentLanguage, "contentLanguage return from CreateVideoEnrollment() is not \"en-US\"")
	assert.Equal("en-US", cver2.ContentLanguage, "contentLanguage return from CreateVideoEnrollment() is not \"en-US\"")
	assert.Equal("en-US", cver3.ContentLanguage, "contentLanguage return from CreateVideoEnrollment() is not \"en-US\"")
	assert.Equal(201, cver1.Status, "status return from CreateVideoEnrollment() is not 201")
	assert.Equal(201, cver2.Status, "status return from CreateVideoEnrollment() is not 201")
	assert.Equal(201, cver3.Status, "status return from CreateVideoEnrollment() is not 201")
	assert.Equal("Never forget tomorrow is a new day", cver1.Text, "text return from CreateVideoEnrollment() from videoEnrollmentArmaan1.mov is not \"never forget tomorrow is a new day\"")
	assert.Equal("Never forget tomorrow is a new day", cver2.Text, "text return from CreateVideoEnrollment() from videoEnrollmentArmaan2.mov is not \"never forget tomorrow is a new day\"")
	assert.Equal("Never forget tomorrow is a new day", cver3.Text, "text return from CreateVideoEnrollment() from videoEnrollmentArmaan3.mov is not \"never forget tomorrow is a new day\"")
	assert.Equal("SUCC", cver1.ResponseCode, "responseCode return from CreateVideoEnrollment() is not \"SUCC\"")
	assert.Equal("SUCC", cver2.ResponseCode, "responseCode return from CreateVideoEnrollment() is not \"SUCC\"")
	assert.Equal("SUCC", cver3.ResponseCode, "responseCode return from CreateVideoEnrollment() is not \"SUCC\"")

	// Verify
	var vvr VideoVerificationReturn
	err5 := json.Unmarshal([]byte(myVoiceIt.VideoVerification(userId, "en-US", "./videoVerificationArmaan1.mov", false)), &vvr)
	if err5 != nil {
		t.Error(err5.Error())
	}
	r2, _ := regexp.Compile("Successfully verified user with userId : usr_([a-z0-9]){32}")
	assert.True(r2.MatchString(vvr.Message), "message return from VideoVerification() does not follow the pattern \"Successfully verified user with userId : usr_00000000000000000000000000000000\"")
	assert.Equal(200, vvr.Status, "status return from VideoVerification() is not 200")
	assert.Equal("SUCC", vvr.ResponseCode, "responseCode return from VideoVerification() is not \"SUCC\"")

	// Identify
	// Create user to add users to
	var cgr CreateGroupReturn
	err6 := json.Unmarshal([]byte(myVoiceIt.CreateGroup("")), &cgr)
	if err6 != nil {
		t.Error(err6.Error())
	}
	groupId := cgr.GroupId
	myVoiceIt.AddUserToGroup(groupId, userId)

	// Create another user to use VideoIdentification()
	var cur2 CreateUserReturn
	cur2err := json.Unmarshal([]byte(myVoiceIt.CreateUser()), &cur2)
	if cur2err != nil {
		t.Error(err1.Error())
	}
	userId2 := cur2.UserId
	myVoiceIt.AddUserToGroup(groupId, userId2)
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoEnrollmentStephen1.mov")
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoEnrollmentStephen2.mov")
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoEnrollmentStephen3.mov")

	var cve4 CreateVideoEnrollmentReturn
	cveerr1 := json.Unmarshal([]byte(myVoiceIt.CreateVideoEnrollment(userId2, "en-US", "./videoEnrollmentStephen1.mov", false)), &cve4)
	if cveerr1 != nil {
		t.Error(cveerr1.Error())
	}
	var cve5 CreateVideoEnrollmentReturn
	cveerr2 := json.Unmarshal([]byte(myVoiceIt.CreateVideoEnrollment(userId2, "en-US", "./videoEnrollmentStephen2.mov", false)), &cve5)
	if cveerr2 != nil {
		t.Error(cveerr2.Error())
	}
	var cve6 CreateVideoEnrollmentReturn
	cveerr3 := json.Unmarshal([]byte(myVoiceIt.CreateVideoEnrollment(userId2, "en-US", "./videoEnrollmentStephen3.mov", false)), &cve6)
	if cveerr3 != nil {
		t.Error(cveerr3.Error())
	}

	// Video Identification
	var vir VideoIdentificationReturn
	err7 := json.Unmarshal([]byte(myVoiceIt.VideoIdentification(groupId, "en-US", "./videoVerificationArmaan1.mov")), &vir)
	if err7 != nil {
		t.Error(err7.Error())
	}

	assert.Equal(userId, vir.UserId, "VideoIdentification() failed to identify user "+userId+" from group "+groupId)

	// Delete Enrollments
	var der1 DeleteEnrollmentReturn
	err8 := json.Unmarshal([]byte(myVoiceIt.DeleteEnrollment(userId, strconv.Itoa(cver1.Id))), &der1)
	if err8 != nil {
		t.Error(err8.Error())
	}
	var der2 DeleteEnrollmentReturn
	der2str := myVoiceIt.DeleteEnrollment(userId, strconv.Itoa(cver2.Id))
	err9 := json.Unmarshal([]byte(der2str), &der2)
	if err9 != nil {
		t.Error(err9.Error())
	}

	var der3 DeleteEnrollmentReturn
	err10 := json.Unmarshal([]byte(myVoiceIt.DeleteEnrollment(userId, strconv.Itoa(cver3.Id))), &der3)
	if err10 != nil {
		t.Error(err10.Error())
	}
	r3, _ := regexp.Compile("Deleted enrollment with id : " + strconv.Itoa(cver1.Id))
	assert.True(r3.MatchString(der1.Message), "message return from DeleteEnrollment() does not follow the pattern \"Deleted enrollment with id : 0\"")

	r4, _ := regexp.Compile("Deleted enrollment with id : " + strconv.Itoa(cver2.Id))
	assert.True(r4.MatchString(der2.Message), "message return from DeleteEnrollment() does not follow the pattern \"Deleted enrollment with id : 0\"")

	r5, _ := regexp.Compile("Deleted enrollment with id : " + strconv.Itoa(cver3.Id))
	assert.True(r5.MatchString(der3.Message), "message return from DeleteEnrollment() does not follow the pattern \"Deleted enrollment with id : 0\"")

	myVoiceIt.RemoveUserFromGroup(groupId, userId)
	myVoiceIt.RemoveUserFromGroup(groupId, userId2)
	myVoiceIt.DeleteGroup(groupId)
	myVoiceIt.DeleteUser(userId)
	myVoiceIt.DeleteAllEnrollmentsForUser(userId2)
	myVoiceIt.DeleteUser(userId2)

	os.Remove("./videoEnrollmentArmaan1.mov")
	os.Remove("./videoEnrollmentArmaan2.mov")
	os.Remove("./videoEnrollmentArmaan3.mov")
	os.Remove("./videoVerificationArmaan1.mov")
	os.Remove("./videoEnrollmentStephen1.mov")
	os.Remove("./videoEnrollmentStephen2.mov")
	os.Remove("./videoEnrollmentStephen3.mov")

}

// Test video enrollment/verification/identification by URL (and deleting each individually)
func TestVideoEnrollmentVerificationIdentificationByUrl(t *testing.T) {
	assert := assert.New(t)
	apikey := os.Getenv("VIAPIKEY")
	apitoken := os.Getenv("VIAPITOKEN")
	myVoiceIt := voiceit2.NewClient(apikey, apitoken)

	// CreateUser()
	var cur CreateUserReturn
	err1 := json.Unmarshal([]byte(myVoiceIt.CreateUser()), &cur)
	if err1 != nil {
		t.Error(err1.Error())
	}
	userId := cur.UserId

	// Make 3 enrollments

	// Enrollment1
	var cver1 CreateVideoEnrollmentReturn
	err2 := json.Unmarshal([]byte(myVoiceIt.CreateVideoEnrollmentByUrl(userId, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoEnrollmentArmaan1.mov", false)), &cver1)
	if err2 != nil {
		t.Error(err2.Error())
	}

	// Enrollment2
	var cver2 CreateVideoEnrollmentReturn
	err3 := json.Unmarshal([]byte(myVoiceIt.CreateVideoEnrollmentByUrl(userId, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoEnrollmentArmaan2.mov", false)), &cver2)
	if err3 != nil {
		t.Error(err3.Error())
	}

	// Enrollment3
	var cver3 CreateVideoEnrollmentReturn
	err4 := json.Unmarshal([]byte(myVoiceIt.CreateVideoEnrollmentByUrl(userId, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoEnrollmentArmaan3.mov", false)), &cver3)
	if err4 != nil {
		t.Error(err4.Error())
	}

	// Run checks on enrollment returns
	r1, _ := regexp.Compile("Successfully added video enrollment for user with userId : usr_([a-z0-9]){32}")
	assert.True(r1.MatchString(cver1.Message), "message return from CreateVideoEnrollmentByUrl() does not follow the pattern \"Successfully added video enrollment for user with userId : usr_00000000000000000000000000000000\"")
	assert.True(r1.MatchString(cver2.Message), "message return from CreateVideoEnrollmentByUrl() does not follow the pattern \"Successfully added video enrollment for user with userId : usr_00000000000000000000000000000000\"")
	assert.True(r1.MatchString(cver3.Message), "message return from CreateVideoEnrollmentByUrl() does not follow the pattern \"Successfully added video enrollment for user with userId : usr_00000000000000000000000000000000\"")
	assert.Equal("en-US", cver1.ContentLanguage, "contentLanguage return from CreateVideoEnrollmentByUrl() is not \"en-US\"")
	assert.Equal("en-US", cver2.ContentLanguage, "contentLanguage return from CreateVideoEnrollmentByUrl() is not \"en-US\"")
	assert.Equal("en-US", cver3.ContentLanguage, "contentLanguage return from CreateVideoEnrollmentByUrl() is not \"en-US\"")
	assert.Equal(201, cver1.Status, "status return from CreateVideoEnrollmentByUrl() is not 201")
	assert.Equal(201, cver2.Status, "status return from CreateVideoEnrollmentByUrl() is not 201")
	assert.Equal(201, cver3.Status, "status return from CreateVideoEnrollmentByUrl() is not 201")
	assert.Equal("Never forget tomorrow is a new day", cver1.Text, "text return from CreateVideoEnrollmentByUrl() from videoEnrollmentArmaan1.mov is not \"never forget tomorrow is a new day\"")
	assert.Equal("Never forget tomorrow is a new day", cver2.Text, "text return from CreateVideoEnrollmentByUrl() from videoEnrollmentArmaan2.mov is not \"never forget tomorrow is a new day\"")
	assert.Equal("Never forget tomorrow is a new day", cver3.Text, "text return from CreateVideoEnrollmentByUrl() from videoEnrollmentArmaan3.mov is not \"never forget tomorrow is a new day\"")
	assert.Equal("SUCC", cver1.ResponseCode, "responseCode return from CreateVideoEnrollmentByUrl() is not \"SUCC\"")
	assert.Equal("SUCC", cver2.ResponseCode, "responseCode return from CreateVideoEnrollmentByUrl() is not \"SUCC\"")
	assert.Equal("SUCC", cver3.ResponseCode, "responseCode return from CreateVideoEnrollmentByUrl() is not \"SUCC\"")

	// Verify
	var vvr VideoVerificationReturn
	err5 := json.Unmarshal([]byte(myVoiceIt.VideoVerificationByUrl(userId, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoVerificationArmaan1.mov", false)), &vvr)
	if err5 != nil {
		t.Error(err5.Error())
	}
	r2, _ := regexp.Compile("Successfully verified user with userId : usr_([a-z0-9]){32}")
	assert.True(r2.MatchString(vvr.Message), "message return from VideoVerification() does not follow the pattern \"Successfully verified user with userId : usr_00000000000000000000000000000000\"")
	assert.Equal(200, vvr.Status, "status return from VideoVerification() is not 200")
	assert.Equal("SUCC", vvr.ResponseCode, "responseCode return from VideoVerification() is not \"SUCC\"")

	// Identify
	// Create user to add users to
	var cgr CreateGroupReturn
	err6 := json.Unmarshal([]byte(myVoiceIt.CreateGroup("")), &cgr)
	if err6 != nil {
		t.Error(err6.Error())
	}
	groupId := cgr.GroupId
	myVoiceIt.AddUserToGroup(groupId, userId)

	// Create another user to use VideoIdentification()
	var cur2 CreateUserReturn
	cur2err := json.Unmarshal([]byte(myVoiceIt.CreateUser()), &cur2)
	if cur2err != nil {
		t.Error(err1.Error())
	}
	userId2 := cur2.UserId
	myVoiceIt.AddUserToGroup(groupId, userId2)

	var cve4 CreateVideoEnrollmentReturn
	cveerr1 := json.Unmarshal([]byte(myVoiceIt.CreateVideoEnrollmentByUrl(userId2, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoEnrollmentStephen1.mov", false)), &cve4)
	if cveerr1 != nil {
		t.Error(cveerr1.Error())
	}
	var cve5 CreateVideoEnrollmentReturn
	cveerr2 := json.Unmarshal([]byte(myVoiceIt.CreateVideoEnrollmentByUrl(userId2, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoEnrollmentStephen2.mov", false)), &cve5)
	if cveerr2 != nil {
		t.Error(cveerr2.Error())
	}
	var cve6 CreateVideoEnrollmentReturn
	cveerr3 := json.Unmarshal([]byte(myVoiceIt.CreateVideoEnrollmentByUrl(userId2, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoEnrollmentStephen3.mov", false)), &cve6)
	if cveerr3 != nil {
		t.Error(cveerr3.Error())
	}

	// Video Identification
	var vir VideoIdentificationReturn
	err7 := json.Unmarshal([]byte(myVoiceIt.VideoIdentificationByUrl(groupId, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/videoVerificationArmaan1.mov")), &vir)
	if err7 != nil {
		t.Error(err7.Error())
	}

	assert.Equal(userId, vir.UserId, "VideoIdentificationByUrl() failed to identify user "+userId+" from group "+groupId)

	// Delete Enrollments
	var der1 DeleteEnrollmentReturn
	err8 := json.Unmarshal([]byte(myVoiceIt.DeleteEnrollment(userId, strconv.Itoa(cver1.Id))), &der1)
	if err8 != nil {
		t.Error(err8.Error())
	}
	var der2 DeleteEnrollmentReturn
	der2str := myVoiceIt.DeleteEnrollment(userId, strconv.Itoa(cver2.Id))
	err9 := json.Unmarshal([]byte(der2str), &der2)
	if err9 != nil {
		t.Error(err9.Error())
	}

	var der3 DeleteEnrollmentReturn
	err10 := json.Unmarshal([]byte(myVoiceIt.DeleteEnrollment(userId, strconv.Itoa(cver3.Id))), &der3)
	if err10 != nil {
		t.Error(err10.Error())
	}
	r3, _ := regexp.Compile("Deleted enrollment with id : " + strconv.Itoa(cver1.Id))
	assert.True(r3.MatchString(der1.Message), "message return from DeleteEnrollment() does not follow the pattern \"Deleted enrollment with id : 0\"")

	r4, _ := regexp.Compile("Deleted enrollment with id : " + strconv.Itoa(cver2.Id))
	assert.True(r4.MatchString(der2.Message), "message return from DeleteEnrollment() does not follow the pattern \"Deleted enrollment with id : 0\"")

	r5, _ := regexp.Compile("Deleted enrollment with id : " + strconv.Itoa(cver3.Id))
	assert.True(r5.MatchString(der3.Message), "message return from DeleteEnrollment() does not follow the pattern \"Deleted enrollment with id : 0\"")

	myVoiceIt.RemoveUserFromGroup(groupId, userId)
	myVoiceIt.RemoveUserFromGroup(groupId, userId2)
	myVoiceIt.DeleteGroup(groupId)
	myVoiceIt.DeleteUser(userId)
	myVoiceIt.DeleteAllEnrollmentsForUser(userId2)
	myVoiceIt.DeleteUser(userId2)

}

// Test voice enrollment/verification/identification

type CreateVoiceEnrollmentReturn struct {
	Message         string  `json:"message"`
	ContentLanguage string  `json:"contentLanguage"`
	Id              int     `json:"id"`
	Status          int     `json:"status"`
	Text            string  `json:"text"`
	TextConfidence  float32 `json:"textConfidence"`
	CreatedAt       int     `json:"createdAt"`
	TimeTaken       string  `json:"timeTaken"`
	ResponseCode    string  `json:"responseCode"`
}

type VoiceVerificationReturn struct {
	Message        string  `json:"message"`
	Status         int     `json:"status"`
	Confidence     float32 `json:"confidence"`
	Text           string  `json:"text"`
	TextConfidence float32 `json:"textConfidence"`
	TimeTaken      string  `json:"timeTaken"`
	ResponseCode   string  `json:"responseCode"`
}

type VoiceIdentificationReturn struct {
	Message        string  `json:"message"`
	UserId         string  `json:"userId"`
	GroupId        string  `json:"groupId"`
	Confidence     float32 `json:"confidence"`
	Status         int     `json:"status"`
	Text           string  `json:"text"`
	TextConfidence float32 `json:"textConfidence"`
	TimeTaken      string  `json:"timeTaken"`
	ResponseCode   string  `json:"responseCode"`
}

func TestVoiceEnrollmentVerificationIdentification(t *testing.T) {
	assert := assert.New(t)
	apikey := os.Getenv("VIAPIKEY")
	apitoken := os.Getenv("VIAPITOKEN")
	myVoiceIt := voiceit2.NewClient(apikey, apitoken)

	// CreateUser() * 2
	var cur1 CreateUserReturn
	err1 := json.Unmarshal([]byte(myVoiceIt.CreateUser()), &cur1)
	if err1 != nil {
		t.Error(err1.Error())
	}
	userId1 := cur1.UserId

	var cur2 CreateUserReturn
	err2 := json.Unmarshal([]byte(myVoiceIt.CreateUser()), &cur2)
	if err2 != nil {
		t.Error(err2.Error())
	}
	userId2 := cur2.UserId

	// Enroll Voice * 3 * 2

	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/enrollmentArmaan1.wav")
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/enrollmentArmaan2.wav")
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/enrollmentArmaan3.wav")
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/verificationArmaan1.wav")

	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/enrollmentStephen1.wav")
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/enrollmentStephen2.wav")
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/enrollmentStephen3.wav")

	var cver1 CreateVoiceEnrollmentReturn
	err3 := json.Unmarshal([]byte(myVoiceIt.CreateVoiceEnrollment(userId1, "en-US", "./enrollmentArmaan1.wav")), &cver1)
	if err3 != nil {
		t.Error(err3.Error())
	}

	var cver2 CreateVoiceEnrollmentReturn
	err4 := json.Unmarshal([]byte(myVoiceIt.CreateVoiceEnrollment(userId1, "en-US", "./enrollmentArmaan2.wav")), &cver2)
	if err4 != nil {
		t.Error(err4.Error())
	}

	var cver3 CreateVoiceEnrollmentReturn
	err5 := json.Unmarshal([]byte(myVoiceIt.CreateVoiceEnrollment(userId1, "en-US", "./enrollmentArmaan3.wav")), &cver3)
	if err5 != nil {
		t.Error(err5.Error())
	}

	var cver4 CreateVoiceEnrollmentReturn
	err6 := json.Unmarshal([]byte(myVoiceIt.CreateVoiceEnrollment(userId2, "en-US", "./enrollmentStephen1.wav")), &cver4)
	if err6 != nil {
		t.Error(err6.Error())
	}

	var cver5 CreateVoiceEnrollmentReturn
	err7 := json.Unmarshal([]byte(myVoiceIt.CreateVoiceEnrollment(userId2, "en-US", "./enrollmentStephen2.wav")), &cver5)
	if err7 != nil {
		t.Error(err7.Error())
	}

	var cver6 CreateVoiceEnrollmentReturn
	err8 := json.Unmarshal([]byte(myVoiceIt.CreateVoiceEnrollment(userId2, "en-US", "./enrollmentStephen3.wav")), &cver6)
	if err8 != nil {
		t.Error(err8.Error())
	}

	// Test Enrollment Returns
	r1, _ := regexp.Compile("Successfully enrolled user with userId : usr_([a-z0-9]){32}")
	assert.True(r1.MatchString(cver1.Message), "message return from CreateVoiceEnrollment() does not follow pattern \"Successfully added voice enrollment for user with userId : usr_00000000000000000000000000000000\"")
	assert.True(r1.MatchString(cver2.Message), "message return from CreateVoiceEnrollment() does not follow pattern \"Successfully added voice enrollment for user with userId : usr_00000000000000000000000000000000\"")
	assert.True(r1.MatchString(cver3.Message), "message return from CreateVoiceEnrollment() does not follow pattern \"Successfully added voice enrollment for user with userId : usr_00000000000000000000000000000000\"")
	assert.Equal("en-US", cver1.ContentLanguage, "contentLanguage return from CreateVoiceEnrollment() not \"en-US\"")
	assert.Equal("en-US", cver2.ContentLanguage, "contentLanguage return from CreateVoiceEnrollment() not \"en-US\"")
	assert.Equal("en-US", cver3.ContentLanguage, "contentLanguage return from CreateVoiceEnrollment() not \"en-US\"")
	assert.Equal("Never forget tomorrow is a new day", cver1.Text, "text return from CreateVoiceEnrollment() not \"Never forget tomorrow is a new day\"")
	assert.Equal("Never forget tomorrow is a new day", cver2.Text, "text return from CreateVoiceEnrollment() not \"Never forget tomorrow is a new day\"")
	assert.Equal("Never forget tomorrow is a new day", cver3.Text, "text return from CreateVoiceEnrollment() not \"Never forget tomorrow is a new day\"")
	assert.Equal(201, cver1.Status, "status return from CreateVoiceEnrollment() not 201")
	assert.Equal(201, cver2.Status, "status return from CreateVoiceEnrollment() not 201")
	assert.Equal(201, cver3.Status, "status return from CreateVoiceEnrollment() not 201")
	assert.Equal("SUCC", cver1.ResponseCode, "responseCode return from CreateVoiceEnrollment() not \"SUCC\"")
	assert.Equal("SUCC", cver2.ResponseCode, "responseCode return from CreateVoiceEnrollment() not \"SUCC\"")
	assert.Equal("SUCC", cver3.ResponseCode, "responseCode return from CreateVoiceEnrollment() not \"SUCC\"")

	// Verification
	var vvr VoiceVerificationReturn
	err9 := json.Unmarshal([]byte(myVoiceIt.VoiceVerification(userId1, "en-US", "./verificationArmaan1.wav")), &vvr)
	if err9 != nil {
		t.Error(err9.Error())
	}

	r2, _ := regexp.Compile("Successfully verified user with userId : usr_([a-z0-9]){32}")
	assert.True(r2.MatchString(vvr.Message), "message return from VoiceVerification() does not follow pattern \"Successfully verified user with userId : usr_00000000000000000000000000000000\"")
	assert.Equal(200, vvr.Status, "status return from VoiceVerification() is not 200")
	assert.Equal("SUCC", vvr.ResponseCode, "responseCode return from VoiceVerification() is not \"SUCC\"")

	// Identification
	var cgr CreateGroupReturn
	err10 := json.Unmarshal([]byte(myVoiceIt.CreateGroup("Sample Group Description")), &cgr)
	if err10 != nil {
		t.Error(err10.Error())
	}
	groupId := cgr.GroupId
	myVoiceIt.AddUserToGroup(groupId, userId1)
	myVoiceIt.AddUserToGroup(groupId, userId2)

	var vir VoiceIdentificationReturn
	err11 := json.Unmarshal([]byte(myVoiceIt.VoiceIdentification(groupId, "en-US", "./verificationArmaan1.wav")), &vir)
	if err11 != nil {
		t.Error(err11.Error())
	}

	r3, _ := regexp.Compile("Successfully identified user with userId : usr_([a-z0-9]+){32} in group with groupId : grp_([a-z0-9]){32}")
	assert.True(r3.MatchString(vir.Message), "message return from VoiceIdentification() does not follow the pattern \"Successfully identified user with userId : usr_00000000000000000000000000000000 in group with groupId : grp_00000000000000000000000000000000\"")
	assert.Equal(userId1, vir.UserId, "userId return from VoiceIdentification() is different from true userId")
	assert.Equal(200, vir.Status, "status return from VoiceIdentification() not 200")
	assert.Equal("SUCC", vir.ResponseCode, "responseCode return from VoiceIdentification() not \"SUCC\"")

	// Clean Up
	myVoiceIt.RemoveUserFromGroup(groupId, userId1)
	myVoiceIt.RemoveUserFromGroup(groupId, userId2)
	myVoiceIt.DeleteAllEnrollmentsForUser(userId1)
	myVoiceIt.DeleteAllEnrollmentsForUser(userId2)
	myVoiceIt.DeleteGroup(groupId)
	myVoiceIt.DeleteUser(userId1)
	myVoiceIt.DeleteUser(userId2)

	os.Remove("enrollmentArmaan1.wav")
	os.Remove("enrollmentArmaan2.wav")
	os.Remove("enrollmentArmaan3.wav")
	os.Remove("verificationArmaan1.wav")
	os.Remove("enrollmentStephen1.wav")
	os.Remove("enrollmentStephen2.wav")
	os.Remove("enrollmentStephen3.wav")

}

func TestVoiceEnrollmentVerificationIdentificationByUrl(t *testing.T) {
	assert := assert.New(t)
	apikey := os.Getenv("VIAPIKEY")
	apitoken := os.Getenv("VIAPITOKEN")
	myVoiceIt := voiceit2.NewClient(apikey, apitoken)

	// CreateUser() * 2
	var cur1 CreateUserReturn
	err1 := json.Unmarshal([]byte(myVoiceIt.CreateUser()), &cur1)
	if err1 != nil {
		t.Error(err1.Error())
	}
	userId1 := cur1.UserId

	var cur2 CreateUserReturn
	err2 := json.Unmarshal([]byte(myVoiceIt.CreateUser()), &cur2)
	if err2 != nil {
		t.Error(err2.Error())
	}
	userId2 := cur2.UserId

	// Enroll Voice * 3 * 2

	var cver1 CreateVoiceEnrollmentReturn
	err3 := json.Unmarshal([]byte(myVoiceIt.CreateVoiceEnrollmentByUrl(userId1, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/enrollmentArmaan1.wav")), &cver1)
	if err3 != nil {
		t.Error(err3.Error())
	}

	var cver2 CreateVoiceEnrollmentReturn
	err4 := json.Unmarshal([]byte(myVoiceIt.CreateVoiceEnrollmentByUrl(userId1, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/enrollmentArmaan2.wav")), &cver2)
	if err4 != nil {
		t.Error(err4.Error())
	}

	var cver3 CreateVoiceEnrollmentReturn
	err5 := json.Unmarshal([]byte(myVoiceIt.CreateVoiceEnrollmentByUrl(userId1, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/enrollmentArmaan3.wav")), &cver3)
	if err5 != nil {
		t.Error(err5.Error())
	}

	var cver4 CreateVoiceEnrollmentReturn
	err6 := json.Unmarshal([]byte(myVoiceIt.CreateVoiceEnrollmentByUrl(userId2, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/enrollmentStephen1.wav")), &cver4)
	if err6 != nil {
		t.Error(err6.Error())
	}

	var cver5 CreateVoiceEnrollmentReturn
	err7 := json.Unmarshal([]byte(myVoiceIt.CreateVoiceEnrollmentByUrl(userId2, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/enrollmentStephen2.wav")), &cver5)
	if err7 != nil {
		t.Error(err7.Error())
	}

	var cver6 CreateVoiceEnrollmentReturn
	err8 := json.Unmarshal([]byte(myVoiceIt.CreateVoiceEnrollmentByUrl(userId2, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/enrollmentStephen3.wav")), &cver6)
	if err8 != nil {
		t.Error(err8.Error())
	}

	// Test Enrollment Returns
	r1, _ := regexp.Compile("Successfully enrolled user with userId : usr_([a-z0-9]){32}")
	assert.True(r1.MatchString(cver1.Message), "message return from CreateVoiceEnrollmentByUrl() does not follow pattern \"Successfully added voice enrollment for user with userId : usr_00000000000000000000000000000000\"")
	assert.True(r1.MatchString(cver2.Message), "message return from CreateVoiceEnrollmentByUrl() does not follow pattern \"Successfully added voice enrollment for user with userId : usr_00000000000000000000000000000000\"")
	assert.True(r1.MatchString(cver3.Message), "message return from CreateVoiceEnrollmentByUrl() does not follow pattern \"Successfully added voice enrollment for user with userId : usr_00000000000000000000000000000000\"")
	assert.Equal("en-US", cver1.ContentLanguage, "contentLanguage return from CreateVoiceEnrollmenByUrlt() not \"en-US\"")
	assert.Equal("en-US", cver2.ContentLanguage, "contentLanguage return from CreateVoiceEnrollmentByUrl() not \"en-US\"")
	assert.Equal("en-US", cver3.ContentLanguage, "contentLanguage return from CreateVoiceEnrollmentByUrl() not \"en-US\"")
	assert.Equal("Never forget tomorrow is a new day", cver1.Text, "text return from CreateVoiceEnrollmenByUrlt() not \"Never forget tomorrow is a new day\"")
	assert.Equal("Never forget tomorrow is a new day", cver2.Text, "text return from CreateVoiceEnrollmenByUrlt() not \"Never forget tomorrow is a new day\"")
	assert.Equal("Never forget tomorrow is a new day", cver3.Text, "text return from CreateVoiceEnrollmentByUrl() not \"Never forget tomorrow is a new day\"")
	assert.Equal(201, cver1.Status, "status return from CreateVoiceEnrollmentByUrl() not 201")
	assert.Equal(201, cver2.Status, "status return from CreateVoiceEnrollmentByUrl() not 201")
	assert.Equal(201, cver3.Status, "status return from CreateVoiceEnrollmentByUrl() not 201")
	assert.Equal("SUCC", cver1.ResponseCode, "responseCode return from CreateVoiceEnrollmenByUrlt() not \"SUCC\"")
	assert.Equal("SUCC", cver2.ResponseCode, "responseCode return from CreateVoiceEnrollmentByUrl() not \"SUCC\"")
	assert.Equal("SUCC", cver3.ResponseCode, "responseCode return from CreateVoiceEnrollmentByUrl() not \"SUCC\"")

	// Verification By Url
	var vvr VoiceVerificationReturn
	err9 := json.Unmarshal([]byte(myVoiceIt.VoiceVerificationByUrl(userId1, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/verificationArmaan1.wav")), &vvr)
	if err9 != nil {
		t.Error(err9.Error())
	}

	r2, _ := regexp.Compile("Successfully verified user with userId : usr_([a-z0-9]){32}")
	assert.True(r2.MatchString(vvr.Message), "message return from VoiceVerificationByUrl() does not follow pattern \"Successfully verified user with userId : usr_00000000000000000000000000000000\"")
	assert.Equal(200, vvr.Status, "status return from VoiceVerificationByUrl() is not 200")
	assert.Equal("SUCC", vvr.ResponseCode, "responseCode return from VoiceVerificationByUrl() is not \"SUCC\"")

	// Identification By Url
	var cgr CreateGroupReturn
	err10 := json.Unmarshal([]byte(myVoiceIt.CreateGroup("Sample Group Description")), &cgr)
	if err10 != nil {
		t.Error(err10.Error())
	}
	groupId := cgr.GroupId
	myVoiceIt.AddUserToGroup(groupId, userId1)
	myVoiceIt.AddUserToGroup(groupId, userId2)

	var vir VoiceIdentificationReturn
	err11 := json.Unmarshal([]byte(myVoiceIt.VoiceIdentificationByUrl(groupId, "en-US", "https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/verificationArmaan1.wav")), &vir)
	if err11 != nil {
		t.Error(err11.Error())
	}

	r3, _ := regexp.Compile("Successfully identified user with userId : usr_([a-z0-9]+){32} in group with groupId : grp_([a-z0-9]){32}")
	assert.True(r3.MatchString(vir.Message), "message return from VoiceIdentificationByUrl() does not follow the pattern \"Successfully identified user with userId : usr_00000000000000000000000000000000 in group with groupId : grp_00000000000000000000000000000000\"")
	assert.Equal(userId1, vir.UserId, "userId return from VoiceIdentificationByUrl() is different from true userId")
	assert.Equal(200, vir.Status, "status return from VoiceIdentificationByUrl() not 200")
	assert.Equal("SUCC", vir.ResponseCode, "responseCode return from VoiceIdentificatioByUrln() not \"SUCC\"")

	// Clean Up
	myVoiceIt.RemoveUserFromGroup(groupId, userId1)
	myVoiceIt.RemoveUserFromGroup(groupId, userId2)
	myVoiceIt.DeleteAllEnrollmentsForUser(userId1)
	myVoiceIt.DeleteAllEnrollmentsForUser(userId2)
	myVoiceIt.DeleteGroup(groupId)
	myVoiceIt.DeleteUser(userId1)
	myVoiceIt.DeleteUser(userId2)
}

// Test face enrollment/verification/identification
type CreateFaceEnrollmentReturn struct {
	Message          string `json:"message"`
	Status           int    `json:"status"`
	BlinksCount      int    `json:"blinksCount"`
	CreatedAt        int    `json:"createdAt"`
	TimeTaken        string `json:"timeTaken"`
	FaceEnrollmentId int    `json:"faceEnrollmentId"`
	ResponseCode     string `json:"responseCode"`
}

type FaceVerificationReturn struct {
	Message        string  `json:"message"`
	Status         int     `json:"status"`
	FaceConfidence float32 `json:"faceConfidence"`
	BlinksCount    int     `json:"blinksCount"`
	TimeTaken      string  `json:"timeTaken"`
	ResponseCode   string  `json:"responseCode"`
}

func TestFaceEnrollmentVerificationIdentification(t *testing.T) {
	assert := assert.New(t)
	apikey := os.Getenv("VIAPIKEY")
	apitoken := os.Getenv("VIAPITOKEN")
	myVoiceIt := voiceit2.NewClient(apikey, apitoken)

	// CreateUser() * 2
	var cur CreateUserReturn
	err1 := json.Unmarshal([]byte(myVoiceIt.CreateUser()), &cur)
	if err1 != nil {
		t.Error(err1.Error())
	}
	userId := cur.UserId

	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/faceEnrollmentArmaan1.mp4")
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/faceEnrollmentArmaan2.mp4")
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/faceEnrollmentArmaan3.mp4")
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/faceVerificationArmaan1.mp4")

	var cfer1 CreateFaceEnrollmentReturn
	err2 := json.Unmarshal([]byte(myVoiceIt.CreateFaceEnrollment(userId, "./faceEnrollmentArmaan1.mp4", false)), &cfer1)
	if err2 != nil {
		t.Error(err2.Error())
	}

	var cfer2 CreateFaceEnrollmentReturn
	err3 := json.Unmarshal([]byte(myVoiceIt.CreateFaceEnrollment(userId, "./faceEnrollmentArmaan2.mp4", false)), &cfer2)
	if err3 != nil {
		t.Error(err3.Error())
	}

	var cfer3 CreateFaceEnrollmentReturn
	err4 := json.Unmarshal([]byte(myVoiceIt.CreateFaceEnrollment(userId, "./faceEnrollmentArmaan3.mp4", false)), &cfer3)
	if err4 != nil {
		t.Error(err4.Error())
	}

	r1, _ := regexp.Compile("Successfully enrolled face for user with userId : usr_([a-z0-9]+){32}")
	assert.True(r1.MatchString(cfer1.Message), "message return from CreateFaceEnrollment() does not follow the pattern \"Successfully enrolled face for user with userId : usr_00000000000000000000000000000000\"")
	assert.True(r1.MatchString(cfer2.Message), "message return from CreateFaceEnrollment() does not follow the pattern \"Successfully enrolled face for user with userId : usr_00000000000000000000000000000000\"")
	assert.True(r1.MatchString(cfer3.Message), "message return from CreateFaceEnrollment() does not follow the pattern \"Successfully enrolled face for user with userId : usr_00000000000000000000000000000000\"")
	assert.Equal(201, cfer1.Status, "status return from CreateFaceEnrollment() is not 201")
	assert.Equal(201, cfer2.Status, "status return from CreateFaceEnrollment() is not 201")
	assert.Equal(201, cfer3.Status, "status return from CreateFaceEnrollment() is not 201")
	assert.Equal("SUCC", cfer1.ResponseCode, "responseCode return from CreateFaceEnrollment() is not \"SUCC\"")
	assert.Equal("SUCC", cfer2.ResponseCode, "responseCode return from CreateFaceEnrollment() is not \"SUCC\"")
	assert.Equal("SUCC", cfer3.ResponseCode, "responseCode return from CreateFaceEnrollment() is not \"SUCC\"")

	var fvr FaceVerificationReturn
	err5 := json.Unmarshal([]byte(myVoiceIt.FaceVerification(userId, "./faceVerificationArmaan1.mp4")), &fvr)
	if err5 != nil {
		t.Error(err5.Error())
	}

	r2, _ := regexp.Compile("Successfully verified user with userId : usr_([a-z0-9]+){32}")
	assert.True(r2.MatchString(fvr.Message), "message return from FaceVerification() does not follow the pattern \"Successfully verified user with userId : usr_00000000000000000000000000000000\"")
	assert.Equal(200, fvr.Status, "status return from FaceVerification() is not 200")
	assert.Equal("SUCC", fvr.ResponseCode, "status return from FaceVerification() is not \"SUCC\"")

	myVoiceIt.DeleteAllEnrollmentsForUser(userId)
	myVoiceIt.DeleteUser(userId)

	os.Remove("faceEnrollmentArmaan1.mp4")
	os.Remove("faceEnrollmentArmaan1.mp4")
	os.Remove("faceEnrollmentArmaan1.mp4")
	os.Remove("faceVerificationArmaan1.mp4")
}
