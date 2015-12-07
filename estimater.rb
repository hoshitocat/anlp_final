#! /usr/bin/env ruby

require 'pp'

class Estimater

  DEFAULT_PATH = './train_datas/neko.num'

  def initialize
  end
end

pp File.open(Estimater::DEFAULT_PATH, 'r').each { |content| pp content }
