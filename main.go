struct Instance {
  id string,
  usage bool,
  send_notif bool
}
func main() {
  // create ec2 session
  // create ec2 client

  // create channel of Instance

  // go for each dc, go get list of instances

  // when instance in channel check usage
  // if usage not enough, Instance.usage=false
  // if usage==false and send_notif==true, stop instance

  // if usage==false, send email to prevent the kill
  // save send_notif with instance

}
