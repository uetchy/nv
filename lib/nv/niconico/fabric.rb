module Niconico
  class Fabric
    attr_reader :agent

    def initialize(agent=nil)
      @agent = agent || Mechanize.new
      @agent.verify_mode = OpenSSL::SSL::VERIFY_NONE
    end

    def sign_in(email, password)
      @agent.post(
        "https://secure.nicovideo.jp/secure/login?site=niconico",
        "mail" => email,
        "password" => password
      )
      return self
    end

    def signed_in?
      @agent.cookies.any? {|c| c.name == 'user_session'}
    end
  end
end
