# Demo Instructions: Traditional CI with GitHub Actions and Dagger Pipelines

## Preparation Checklist
- Ensure Docker Engine is installed.
- Install Dagger SDK.
- Install Go or Python and Git.
- Set up GitHub/GitLab runners in advance if needed.

---

## Demo: GitHub Actions

### Clone the Repository

```
git clone https://github.com/iximiuz/todolist.git
# Demo Instructions: Traditional CI with GitHub Actions and Dagger Pipelines
```


## Preparation Checklist
- Ensure Docker Engine is installed.
- Install Dagger SDK.
- Install Go or Python and Git.
- Set up GitHub/GitLab runners in advance if needed.

---

## Demo: GitHub Actions

### Clone the Repository

```
git clone https://github.com/iximiuz/todolist.git
```

### Execute Commands to Check Makefile Behavior 
```
make lint
make test
make up
make test-e2e
make build
make docker-build
make docker-push
```

### If you wish to test functionality 
```
#Create TODO items:
curl -X POST -d '{"task": "Finish the course"}' http://localhost:8080/todos
curl -X POST -d '{"task": "Nail the sales!!!"}' http://localhost:8080/todos

#Retrieve the list of tasks:
curl -X GET http://localhost:8080/todos

#Delete a task:
curl -X DELETE http://localhost:8080/todos/{id}
```

---

## Demo - Dagger Pipelines

### Initiliase Dagger SDK 
```
pwd - /Users/admin/poc-dagger-vs-traditional-ci

$dagger init --sdk=go --source=./dagger
```

### Run the Dagger Functions 
```
dagger call lint --directory-arg .
dagger call test --directory-arg .
dagger call build --source=.
dagger call publish-image --source=.


```