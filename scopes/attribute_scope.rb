require_relative '../scope'

module Tritium
  class Scope::Attribute < Scope::Base
    def initialize(attribute, root, parent)
      @object = @attribute = attribute
      super
    end

    def remove
      @attribute.remove
    end

    def name(value =  nil, &block)
      if value
        @attribute.name = value
      end
      if block
        @attribute.name = open_text_scope_with(@attribute.name, &block)
      end
    end

    def value(value = nil, &block)
      if value
        @attribute.value = value
      end
      if block
        @attribute.value = open_text_scope_with(@attribute.value, &block)
      end
    end
  end
end