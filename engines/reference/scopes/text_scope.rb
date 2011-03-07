require_relative '../scope'

module Tritium::Engines::Reference::Scope
  class Text < Base
    attr :text

    # @private
    def initialize(text, root = nil, parent = nil)
      @text = text
      @root = root || self
      super
    end

    def set(text)
      @text = text
    end

    def replace(matcher, replacement = "", &block)
      if matcher.is_a? String
        matcher = Regexp.new(matcher)
      end

      @text.gsub!(matcher) do |match|
        $~.captures.each_with_index do |arg, index|
          var("#{index + 1}", arg)
        end
        replace = replacement
        (replace = open_text_scope_with(replace, &block)) if block
        replace.gsub(/\\([\d])/) do |var_match|
          var($1)
        end
      end 
    end
  
    def remove(matcher = nil)
      if matcher
        replace(matcher, "")
      else
        @text = ""
      end
    end
    
    def doc(type = 'xhtml', &block)
      parser_klass = Tritium::Engines.xml_parsers[type]
      doc = parser_klass.parse(@text)
      node_scope = Tritium::Engines::Reference::Scope::Node.new(doc, @root, self)
      node_scope.instance_eval(&block)
      @text = node_scope.node.send("to_" + type)
    end
  
    def append(text)
      @text << text
    end

    def prepend(text)
      @text.insert(0,text)
    end

    def rewrite(what)
      replace(env["rewrite_#{what}_matcher"], env["rewrite_#{what}_to"] + env["proxy_domain"])
    end
  end
end
