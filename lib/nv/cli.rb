require 'thor'

class Nv::CLI < Thor
  include ::Niconico::Helper

  desc "dl URL", "Download video"
  def dl(ptr, output=".")
    config = Nv::Config.new(Nv::CONFIG_PATH)
    config.verify_for_authentication!('dl')

    nico = Niconico::Base.new.sign_in(config.email, config.password)

    if mylist?(ptr)
      mylist = nico.mylist(ptr)

      puts "Title : #{mylist.title}"
      puts "Desc  : #{mylist.description}"

      mylist.items.each do |item|
        dl(item.link, output)
      end
    else
      video = nico.video(ptr)

      # Inspect
      puts "Downloading... #{video.title}"

      # Donwload
      video.download

      puts "+ done"
    end
  end

  desc "info URL", "Show video/mylist info"
  def info(ptr)
    config = Nv::Config.new(Nv::CONFIG_PATH)
    config.verify_for_authentication!('info')

    nico = Niconico::Base.new.sign_in(config.email, config.password)

    if mylist?(ptr)
      mylist = nico.mylist(ptr)

      puts "Title : #{mylist.title}"
      puts "Desc  : #{mylist.description}"

      mylist.items.each_with_index do |item, i|
        puts "   #{i+1}. #{item.title}"
      end
    else
      video = nico.video(ptr)

      puts video.title
      puts "=" * 40
      puts video.description
      puts "=" * 40
      puts "URL: #{video.watch_url}"
    end
  end

  desc "config KEY VALUE", "Set config"
  def config(key=nil, value=nil)
    config = Nv::Config.new(Nv::CONFIG_PATH)

    unless key
      puts "config:"
      config.to_h.each do |k, v|
        puts "   #{k}=#{v}"
      end
      return
    end

    if value
      config[key] = value
      config.save
    end

    puts "config: #{key}=#{config[key]}"
  end
end
