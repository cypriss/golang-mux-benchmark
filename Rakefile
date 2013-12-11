require "set"
require "erb"

class BenchmarkParser

  attr_reader :benchmarks, :frameworks, :results

  def initialize data
    @frameworks = Set.new
    @benchmarks = Set.new
    @results    = {}
    @data = data
  end

  def self.parse data
    parser = new(data)
    parser.parse
    parser
  end

  def parse
    @data.split(/\n/).each do |line|
      if matches = line.match(/Benchmark([^_]+)_([^\s]+)\s+(\d+)\s+(\d+)\s+/)
        framework  = matches[1]
        benchmark  = matches[2]
        loop_count = matches[3]
        loop_time  = matches[4]

        @frameworks.add framework
        @benchmarks.add benchmark

        @results[benchmark] ||= {}
        @results[benchmark][framework] = loop_time
      end
    end
  end
end

class ChartsGenerator
  def initialize benchmarks, frameworks, results
    @benchmarks = benchmarks
    @frameworks = frameworks
    @results    = results
  end

  def generate tpl_filename, out_filename
    File.open tpl_filename, "r" do |tpl_file|
      File.open out_filename, "w" do |out_file|
        tpl = ERB.new tpl_file.read
        out_file.write tpl.result(binding)
      end
    end
  end
end

task :default do
  cmd = "go test -bench=. 2>/dev/null"
  puts "running benchmarks..."
  puts cmd
  output = `#{cmd}`

  bp = BenchmarkParser.parse output
  generator = ChartsGenerator.new bp.benchmarks, bp.frameworks, bp.results

  out_filename = "index.html"
  puts "Generating charts..."
  generator.generate "index.erb", out_filename
  puts "charts written #{out_filename}"
end

