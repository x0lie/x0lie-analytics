resource "aws_ecs_cluster" "main" {
  name = "x0lie-analytics"
}

resource "aws_ecs_task_definition" "app" {
  family                   = "x0lie-analytics"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.execution.arn
  task_role_arn            = aws_iam_role.task.arn

  container_definitions = jsonencode([{
    name  = "x0lie-analytics"
    image = "${aws_ecr_repository.main.repository_url}:${var.image_tag}"

    portMappings = [{
      containerPort = 8080
    }]

    # TODO: migrate to AWS Secret Manager or SSM Parameter Store to avoid plaintext in task definition
    environment = [{
      name  = "DATABASE_URL"
      value = "postgres://analytics:${var.db_password}@${aws_db_instance.main.address}/analytics?sslmode=require"
    }]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = "/ecs/x0lie-analytics"
        "awslogs-region"        = "us-east-1"
        "awslogs-stream-prefix" = "ecs"
        "awslogs-create-group"  = "true"
      }
    }
  }])
}

resource "aws_ecs_service" "app" {
  name            = "x0lie-analytics"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.app.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [aws_subnet.public_a.id]
    security_groups  = [aws_security_group.ecs.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.app.arn
    container_name   = "x0lie-analytics"
    container_port   = 8080
  }
}
