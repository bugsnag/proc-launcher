Given("I build the extension executable") do
  Dir.chdir(BUILD_DIR) do
    `go build ../features/fixtures/extension.go`
    expect(File.exists? "extension").to be_truthy
  end
end

When("I run the extension executable with {string}") do |args|
  Dir.chdir(BUILD_DIR) do
    start_process("./extension", args)
  end
end

When("I run the extension executable with arguments:") do |table|
  Dir.chdir(BUILD_DIR) do
    args = table.raw.flatten
    start_process(["./extension"] + args)
  end
end

Then("{string} is present in the standard output stream") do |contents|
  expect(PROCESSES[-1][:stdout].read).to include contents.gsub("\\n", "\n")
end

Then("{string} is present in the standard error stream") do |contents|
  expect(PROCESSES[-1][:stderr].read).to include contents.gsub("\\n", "\n")
end

Given("I build the executable") do
  Dir.chdir(BUILD_DIR) do
    `go build ../main.go`
    expect(File.exists? "main").to be_truthy
  end
end

When("I run the executable with arguments:") do |table|
  Dir.chdir(BUILD_DIR) do
    args = table.raw.flatten
    start_process(["./main"] + args)
    sleep 0.2 # give it a sec to start. not ideal.
    ppid = PROCESSES[-1][:thread][:pid]
    child_cmd = args.first
    output = `ps -x -o pid,ppid,comm | grep #{ppid} | grep #{child_cmd}`
    @child_pid = output.split(' ').first
  end
end

When(/^I send SIG(\w+) to the executable$/) do |signal|
  Process.kill signal, PROCESSES[-1][:thread][:pid]
end

Then(/^the process exited with signal SIG(\w+)$/) do |signal|
  # Check that the launcher executable died as expected
  status = PROCESSES[-1][:thread].value
  expect(status.exited?).to be_truthy
  # This is a check that the *child* process is dead
  expect(@child_pid).not_to be_nil # check if it lived long enough to be captured
  output = `ps -x -o pid,comm | grep #{@child_pid}`
  expect(output).to be_eql ""
end
