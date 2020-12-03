if ARGV.length != 1
    puts "Missing input file"
    exit;
end
 
filename = ARGV[0]
puts "Loading input: '#{filename}'"
fh = open filename
content = fh.read
fh.close

entries = content.split("\n")

class Policy
  attr_reader :char, :arg_1, :arg_2
  def initialize(char, arg_1, arg_2)
    @char = char
    @arg_1 = arg_1.to_i
    @arg_2 = arg_2.to_i
  end
end

class PositionPolicy < Policy
   def password_valid?(password)
    "#{password[arg_1 - 1]}#{password[arg_2-1]}".count(char) == 1
  end
end

class RepeatedCharsPolicy < Policy
  def password_valid?(password)
    char_count = password.count(char)
    char_count <= arg_2 && arg_1 <= char_count
  end
end

puts "# Problem 1"
valid_entries = entries.filter do |entry|
  policy_params, letter, password  = entry.split
  RepeatedCharsPolicy.new(letter[0], policy_params.split('-')[0], policy_params.split('-')[1]).password_valid?(password)
end
puts "#{valid_entries.count} valid passwords out of #{entries.count}"

puts "# Problem  2"
valid_entries = entries.filter do |entry|
  policy_params, letter, password  = entry.split
  PositionPolicy.new(letter[0], policy_params.split('-')[0], policy_params.split('-')[1]).password_valid?(password)
end
puts "#{valid_entries.count} valid passwords out of #{entries.count}"

