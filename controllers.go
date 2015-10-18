package main

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
	"github.com/samalba/dockerclient"
)

// Register contains all information that should
// be passed into a valid POST /register call
type Register struct {
	Name     string `valid:"alphanum,required"`
	Password string `valid:"alphanum,required"`
}

// RegisterSuccess is used upon a successful
// POST /register call
type RegisterSuccess struct {
	Success bool
	Message string
}

// RegisterController handles account creation
func RegisterController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var s Register
	var sr = RegisterSuccess{false, ""}

	json.NewDecoder(r.Body).Decode(&s)
	r.Body.Close()
	// Result ignored since a nil err value tells us
	// the same thing that a result of true does.
	if _, err := govalidator.ValidateStruct(s); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}

	if err := (User{
		Username: s.Name,
		Password: s.Password,
	}.Create()); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}

	sr.Message = "Users has been created"
	sr.Success = true
	w.Write(Marshal(sr))
	return
}

// Login although the same as the Register struct using the same
// struct for both did not feel right and makes the code
// slightly confusing.
type Login struct {
	Name     string `valid:"alphanum,required"`
	Password string `valid:"alphanum,required"`
}

// LoginSuccess is used upon a successful
// POST /login call
type LoginSuccess struct {
	Success bool
	Message string
	Token   string
}

// LoginController handles authentication and token generation
func LoginController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var s Login
	var sr = LoginSuccess{false, "", ""}

	json.NewDecoder(r.Body).Decode(&s)
	r.Body.Close()
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}

	user := User{}
	if err := (user.Polulate(s.Name)); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}

	if err := validPassword(s.Password, user.Password); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(401)
		w.Write(Marshal(sr))
		return
	}

	if err := user.GenerateToken(); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}

	sr.Token = user.Token
	sr.Success = true
	w.Write(Marshal(sr))
	return
}

// ListResourceSuccess is used upon a successful
// GET /resource call
type ListResourceSuccess struct {
	Success   bool
	Message   string
	Resources []Resource
}

// ListResourceController brings back a list of all resources
// that are current available to that user
func ListResourceController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sr := ListResourceSuccess{false, "", nil}

	if err := VerifyToken(r.Header.Get("token")); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(401)
		w.Write(Marshal(sr))
		return
	}

	resources, err := Resource{}.FetchAll()
	if err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}

	sr.Success = true
	sr.Resources = resources
	w.Write(Marshal(sr))
	return
}

// CreateResource validates incomming requests
type CreateResource struct {
	Name string `valid:"alphanum,required"`
	File string `valid:"required"`
}

// CreateResourceSuccess ...
type CreateResourceSuccess struct {
	Success bool
	Message string
}

// CreateResourceController creates a new resource that all users now have access to
func CreateResourceController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var s CreateResource
	sr := CreateResourceSuccess{false, ""}

	if err := VerifyToken(r.Header.Get("token")); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(401)
		w.Write(Marshal(sr))
		return
	}

	json.NewDecoder(r.Body).Decode(&s)
	r.Body.Close()
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}

	if err := (Resource{
		Name: s.Name,
		File: s.File,
	}.Create()); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}

	sr.Success = true
	sr.Message = "Your resource is now available for use in new containers"
	w.Write(Marshal(sr))
	return
}

// CreateMachine validates that the input is valid
// to an extent... Validation could be better for
// this controller.
type CreateMachine struct {
	Name     string `valid:"alphanum,required"`
	Username string `valid:"alphanum,required"`
	IP       string `valid:"required"`
}

// CreateMachineSuccess is returned
// when the controller has successfully
// returned
type CreateMachineSuccess struct {
	Success bool
	Message string
}

// CreateMachineController adds a new machine that can be used to launch containers on
func CreateMachineController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var s CreateMachine
	sr := CreateMachineSuccess{false, ""}

	if err := VerifyToken(r.Header.Get("token")); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(401)
		w.Write(Marshal(sr))
		return
	}

	json.NewDecoder(r.Body).Decode(&s)
	r.Body.Close()
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}

	if err := (Machine{
		Name:     s.Name,
		Username: s.Username,
		IP:       s.IP,
	}.Create()); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}

	sr.Success = true
	sr.Message = "Your machine has been added to the pool of machines"
	w.Write(Marshal(sr))
	return
}

// CreateLease validates that a resource name
// Has been specified
type CreateLease struct {
	ResourceName string `valid:"alphanum,required"` // As set in the resources table
	LeaseName    string `valid:"alphanum,required"` // name of container on remote machine
}

// CreateLeaseSuccess ...
type CreateLeaseSuccess struct {
	Success  bool
	Message  string
	Username string
	Machine  string
	Name     string
}

// CreateLeaseController creates a new machine for a user to use
func CreateLeaseController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var s CreateLease
	sr := CreateLeaseSuccess{false, "", "", "", ""}

	if err := VerifyToken(r.Header.Get("token")); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(401)
		w.Write(Marshal(sr))
		return
	}

	json.NewDecoder(r.Body).Decode(&s)
	r.Body.Close()
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}
	// Fetch machine to deploy to
	m := Machine{}
	if err := m.FetchRand(); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}
	// Fetch type of resource to build
	res := Resource{}
	if err := res.Fetch(s.ResourceName); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}
	// Create a tar file in memory that can be used
	// to send to docker to create a new image without
	// ever writing to disk
	tar, err := MakeTarFile(res.File)
	if err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}
	// Connect to remote machine
	docker, err := dockerclient.NewDockerClient("http://"+m.IP+":5555", nil)
	if err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}
	// Build a docker image on remote machine
	buildImageConfig := &dockerclient.BuildImage{
		Context:        tar,
		RepoName:       res.Name,
		SuppressOutput: false,
	}
	// Try to create image on remote docker machine
	reader, err := docker.BuildImage(buildImageConfig)
	if err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}
	// If called before checking for an error and an error occurs you will
	// get a null pointer dereferenced error.
	defer reader.Close()
	// These two commands are super helpful but it is returned as
	// a stream of json... In a better implementation we would want to
	// send this back to the client as it comes in. For now I have just
	// left it in to help with debugging.
	// text, _ := ioutil.ReadAll(reader)
	// fmt.Println(string(text))
	// Create a container
	containerConfig := &dockerclient.ContainerConfig{
		Image:       res.Name,
		Cmd:         []string{"ping", "8.8.8.8"},
		AttachStdin: true,
		Tty:         true}
	containerID, err := docker.CreateContainer(containerConfig, s.LeaseName)
	if err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}

	// Start the container on randomly chosen machine
	hostConfig := &dockerclient.HostConfig{}
	// For some reason this API returns the full 64 character hex String
	// which can't be used when starting the container as it only takes the
	// first 12 characters. Possibly a bad api or exists for legacy reasons.
	err = docker.StartContainer(containerID[0:12], hostConfig)
	if err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}

	// Save Lease in database
	Lease{
		Name:        s.LeaseName,
		Username:    m.Username,
		MachineName: m.IP,
	}.Create()

	sr.Success = true
	sr.Message = "Your lease was successfully fulfilled"
	sr.Machine = m.IP
	sr.Username = m.Username
	sr.Name = s.LeaseName
	w.Write(Marshal(sr))
	return
}

// ListLeaseSuccess ...
type ListLeaseSuccess struct {
	Success bool
	Message string
	Leases  []Lease
}

// ListLeasesController returns all leases for a specific client
func ListLeasesController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sr := ListLeaseSuccess{false, "", nil}

	if err := VerifyToken(r.Header.Get("token")); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(401)
		w.Write(Marshal(sr))
		return
	}

	leases, err := Lease{}.FetchAll()
	if err != nil {
		sr.Message = err.Error()
		w.WriteHeader(401)
		w.Write(Marshal(sr))
		return
	}

	sr.Success = true
	sr.Leases = leases
	w.Write(Marshal(sr))
	return
}

// DeleteLease requires the name of the
// container so that it can be deleted
type DeleteLease struct {
	Name string `valid:"alphanum,required"`
}

// DeleteLeaseSuccess ...
type DeleteLeaseSuccess struct {
	Success bool
	Message string
}

// DeleteLeaseController removes a lease from the remote machine
func DeleteLeaseController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var s DeleteLease
	sr := DeleteLeaseSuccess{}

	if err := VerifyToken(r.Header.Get("token")); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(401)
		w.Write(Marshal(sr))
		return
	}

	json.NewDecoder(r.Body).Decode(&s)
	r.Body.Close()
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}
	// Get information about lease so we know what machine it is located on
	l := Lease{}
	if err := l.Fetch(s.Name); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}
	// Connect to machine that the lease is on
	docker, err := dockerclient.NewDockerClient("http://"+l.MachineName+":5555", nil)
	if err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}
	// Remove the container
	if err := docker.RemoveContainer(s.Name, true, true); err != nil {
		sr.Message = err.Error()
		w.WriteHeader(400)
		w.Write(Marshal(sr))
		return
	}
	// Remove container from the database
	l.Delete()

	sr.Success = true
	sr.Message = "Your lease has been removed"
	w.Write(Marshal(sr))
	return
}
