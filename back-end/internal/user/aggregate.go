package user

type UserAggregate struct {
	User     User
	Profile  *Profile
	Producer *Producer
}
