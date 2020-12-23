require 'open3'

BUILD_DIR = File.join(Dir.pwd, "build")
PROCESSES = []
VERBOSE = ENV['VERBOSE'] || ARGV.include?('--verbose')

Before do
  FileUtils.mkdir_p BUILD_DIR
  PROCESSES.clear
end

After do
  PROCESSES.each do |p|
    begin
      if VERBOSE
        puts p[:stderr].read
      end
      Process.kill 'KILL', p[:thread][:pid]
    rescue
    end
  end
end

at_exit do
  FileUtils.rm_r BUILD_DIR
end

def start_process args
  stdin, stdout, stderr, thread = Open3.popen3(*args)
  PROCESSES << {
    thread: thread,
    stdout: stdout,
    stderr: stderr,
    stdin: stdin
  }
end
