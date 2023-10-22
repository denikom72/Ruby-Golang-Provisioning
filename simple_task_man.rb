class Task
  attr_accessor :title, :description, :completed

  def initialize(title, description = "")
    @title = title
    @description = description
    @completed = false
  end

  def mark_as_completed
    @completed = true
  end

  def to_s
    if @completed
      "[X] #{@title} - #{@description}"
    else
      "[ ] #{@title} - #{@description}"
    end
  end
end

class TaskList
  def initialize
    @tasks = []
  end

  def add_task(title, description = "")
    task = Task.new(title, description)
    @tasks << task
  end

  def list_tasks
    @tasks.each_with_index do |task, index|
      puts "#{index + 1}. #{task}"
    end
  end

  def mark_task_completed(index)
    task = @tasks[index]
    task.mark_as_completed if task
  end

  def edit_task_description(index, description)
    task = @tasks[index]
    task.description = description if task
  end

  def delete_task(index)
    @tasks.delete_at(index) if @tasks[index]
  end

  def save_tasks_to_file(filename)
    File.open(filename, 'w') do |file|
      @tasks.each do |task|
        file.puts "#{task.title},#{task.description},#{task.completed}"
      end
    end
  end

  def load_tasks_from_file(filename)
    @tasks = []
    File.readlines(filename).each do |line|
      title, description, completed = line.chomp.split(',')
      task = Task.new(title, description)
      task.completed = completed == 'true'
      @tasks << task
    end
  end
end

class TaskManager
  def initialize
    @task_list = TaskList.new
  end

  def run
    loop do
      puts "Task Manager Menu:"
      puts "1. Add Task"
      puts "2. List Tasks"
      puts "3. Mark Task as Completed"
      puts "4. Edit Task Description"
      puts "5. Delete Task"
      puts "6. Save Tasks to File"
      puts "7. Load Tasks from File"
      puts "8. Exit"
      print "Choose an option: "
      choice = gets.to_i

      case choice
      when 1
        print "Enter task title: "
        title = gets.chomp
        print "Enter task description: "
        description = gets.chomp
        @task_list.add_task(title, description)
      when 2
        @task_list.list_tasks
      when 3
        print "Enter the task number to mark as completed: "
        index = gets.to_i - 1
        @task_list.mark_task_completed(index)
      when 4
        print "Enter the task number to edit description: "
        index = gets.to_i - 1
        print "Enter new description: "
        description = gets.chomp
        @task_list.edit_task_description(index, description)
      when 5
        print "Enter the task number to delete: "
        index = gets.to_i - 1
        @task_list.delete_task(index)
      when 6
        print "Enter the filename to save tasks to: "
        filename = gets.chomp
        @task_list.save_tasks_to_file(filename)
      when 7
        print "Enter the filename to load tasks from: "
        filename = gets.chomp
        @task_list.load_tasks_from_file(filename)
      when 8
        break
      else
        puts "Invalid option. Please try again."
      end
    end
  end
end

task_manager = TaskManager.new
task_manager.run
