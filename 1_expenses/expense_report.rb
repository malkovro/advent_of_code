if ARGV.length != 1
    puts "We need exactly one parameter. The name of a file."
    exit;
end
 
filename = ARGV[0]
puts "Going to open '#{filename}'"
 
fh = open filename
 
content = fh.read
 
fh.close

expenses = content.split.map(&:to_i)
n = expenses.count

expenses.each_with_index do |expense, i|
  expenses.slice(i, n-i).each_with_index do |exp_b, j|
    expenses.slice(i+j, n-(i+j)).each do |exp_c|
      if(exp_b + exp_c + expense == 2020)
        p "#{expense} + #{exp_b} + #{exp_c} = 2020"
        p expense * exp_b * exp_c
      end
    end
  end
end

expenses.each_with_index do |expense, i|
  expenses.slice(i, n-i).each do |exp_b|
    if(exp_b + expense == 2020)
      p "#{expense} + #{exp_b} = 2020"
      p expense * exp_b
    end
  end
end
