resource "aws_db_subnet_group" "main" {
  name       = "x0lie-analytics"
  subnet_ids = [aws_subnet.private_a.id, aws_subnet.private_b.id]
}

resource "aws_db_instance" "main" {
  identifier        = "x0lie-analytics"
  engine            = "postgres"
  engine_version    = "17"
  instance_class    = "db.t4g.micro"
  allocated_storage = 20

  db_name  = "analytics"
  username = "analytics"
  password = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  skip_final_snapshot = true
}
