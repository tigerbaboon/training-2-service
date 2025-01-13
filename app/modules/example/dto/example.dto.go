package dto

type ExampleDTOResponse struct {
	json      string `naming:"camel_case"`
	Name      string
	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
	ID        int64
}

type UserDTOResponse struct {
	Name      string
	CreatedAt int64 `copier:"CreatedAt"`
	UpdatedAt int64
	DeletedAt int64
	ID        int64
}

func (usr *UserDTOResponse) FirstName(firstName string) {
	usr.Name = firstName + "test"
}
