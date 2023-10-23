# VoipProvisGUI Class
#
# This class provides methods for generating a textual form from a model, product, and brand.
# It includes functionality to convert structured data into HTML forms.
#
# Class Methods:
# - initialize(base): Initializes the class with a base directory path.
# - generate_textual_form(model, product, brand, hide_array = []): Generates a textual form based on the provided data.
# - generate_complete_array(model, product, brand): Generates a structured data array.
# - arraysearchrecursive(needle, haystack, needle_key = '', strict = false, path = []): Recursive search in an array.
# - convert2html(key, data): Converts structured data into HTML form elements.
# - file2json(file): Reads a JSON file and returns its contents as a Hash.
#
# Example Usage:
#
# Initialize the class:
# voip_provis = VoipProvisGUI.new('/path/to/base_directory')
#
# Generate a textual form:
# form_data = voip_provis.generate_textual_form('model_name', 'product_name', 'brand_name', ['hide_item1', 'hide_item2'])
#
# Author: Denis Komnenovic
# Version: 1.3
class VoipProvisGUI
  def initialize(base)
    @base = base
  end

  # Generates a textual form based on the provided model, product, and brand.
  def generate_textual_form(model, product, brand, hide_array = [])
    template_array = generate_complete_array(model, product, brand)
    html = {}
    template_array['data'].each do |category, subs|
      subs.each do |subcategories, its|
        its.each do |kitems, items|
          unless hide_array.include?(kitems)
            if kitems =~ /^option\|(.*)/i
              html[category] ||= {}
              html[category][subcategories] ||= {}
              html[category][subcategories][kitems] = convert2html(kitems, items[0])
            end

            if kitems =~ /^loop\|(.*)/i
              items.each do |loop_key, loop_data|
                key = "#{kitems}|#{loop_key}"
                html[category] ||= {}
                html[category][subcategories] ||= {}
                html[category][subcategories][key] = convert2html(key, loop_data)
              end
            end

            if kitems =~ /^lineloop\|(.*)/i
              items.each do |loop_key, loop_data|
                split = kitems.split('_')
                line = split[1]
                key = "lineloop|#{line}|#{loop_key}"
                html[category] ||= {}
                html[category][subcategories] ||= {}
                html[category][subcategories][key] = convert2html(key, loop_data)
              end
            end

            if kitems =~ /^break/
              html[category] ||= {}
              html[category][subcategories] ||= []
              html[category][subcategories] << '<br />'
            end
          end
        end
      end
    end
    html
  end

  # Generates a structured data array based on the model, product, and brand.
  def generate_complete_array(model, product, brand)
    data = {}
    fd_json = file2json(File.join(@base, brand, product, 'family_data.json'))
    model_location = arraysearchrecursive(model, fd_json, 'model')

    raise Exception, 'cant find model' unless model_location

    model_information = fd_json['data']['model_list'][model_location[2]]

    data['phone_data'] = {
      'brand' => brand,
      'product' => product,
      'model' => model,
    }
    data['lines'] = model_information['lines']
    files = model_information['template_data']
    files.unshift('/../../global_template_data.json')
    b = 0
    files.each do |files_data|
      file_path = File.join(@base, brand, product, files_data)
      if File.exist?(file_path)
        template_data = file2json(file_path)
        template_data['template_data']['category'].each do |category|
          category_name = category['name']
          category['subcategory'].each do |subcategory|
            subcategory_name = subcategory['name']
            items_fin = []
            items_loop = []
            subcategory['item'].each do |item|
              case item['type']
              when 'loop_line_options'
                (1..data['lines']).each do |i|
                  var_nam = "lineloop|line_#{i}"
                  item['data']['item'].each do |item_loop|
                    next if item_loop['type'] == 'break'
                    z = item_loop['variable'].delete('$')
                    items_loop[var_nam] ||= {}
                    items_loop[var_nam][z] = item_loop
                    items_loop[var_nam][z]['description'].gsub!('{$count}', i.to_s)
                    items_loop[var_nam][z]['default_value'].gsub!('{$count}', i.to s)
                    items_loop[var_nam][z]['line_loop'] = true
                    items_loop[var_nam][z]['line_count'] = i
                  end
                  items_loop[var_nam] ||= [{'type' => 'break'}]
                end
                items_fin.concat(items_loop)
              when 'loop'
                (item['loop_start']..item['loop_end']).each do |i|
                  name = item['data']['item'][0]['variable'].split('_')
                  var_nam = "loop|#{name[0].delete('$')}_#{i}"
                  item['data']['item'].each do |item_loop|
                    next if item_loop['type'] == 'break'
                    z_tmp = item_loop['variable'].split('_')
                    z = z_tmp[1]
                    items_loop[var_nam] ||= {}
                    items_loop[var_nam][z] = item_loop
                    items_loop[var_nam][z]['description'].gsub!('{$count}', i.to_s)
                    item['variable'].gsub!('_', "_#{i}_")
                    item_loop['default_value'] ||= ''
                    items_loop[var_nam][z]['loop'] = true
                    items_loop[var_nam][z]['loop_count'] = i
                  end
                  items_fin.concat(items_loop)
                end
              when 'break'
                items_fin << 'break'
              else
                var_nam = "option|#{item['variable'].delete('$')}"
                items_fin[var_nam] ||= []
                items_fin[var_nam] << item
              end
            end
            if data['data'][category_name] && data['data'][category_name][subcategory_name]
              old_sc = data['data'][category_name][subcategory_name]
              sub_cat_data = {}
              sub_cat_data[category_name] ||= {}
              sub_cat_data[category_name][subcategory_name] = old_sc + items_fin
              data['data'][category_name] = old_c + new_c
            else
              sub_cat_data[category_name] ||= {}
              sub_cat_data[category_name][subcategory_name] = items_fin
            end
          end
          if data['data'][category_name]
            old_c = data['data'][category_name]
            new_c = sub_cat_data[category_name]
            data['data'][category_name] = old_c + new_c
          else
            data['data'][category_name] = sub_cat_data[category_name]
          end
        end
      end
    end
    data
  end

  # Recursive search in an array.
  def arraysearchrecursive(needle, haystack, needle_key = '', strict = false, path = [])
    return false unless haystack.is_a?(Array)

    haystack.each_with_index do |val, key|
      if val.is_a?(Array) && sub_path = arraysearchrecursive(needle, val, needle_key, strict, path)
        path = path.concat([key]).concat(sub_path)
        return path
      elsif (!strict && val == needle && key == (needle_key.length > 0 ? needle_key : key)) ||
            (strict && val === needle && key == (needle_key.length > 0 ? needle_key : key))
        path << key
        return path
      end
    end
    false
  end

  # Converts structured data into HTML form elements.
  def convert2html(key, data)
    html_return = ''
    case data['type']
    when 'input'
      value = data['value'].to_s.empty? ? data['default_value'] : data['value']
      html_return = "#{data['description']}: <input type='text' name='#{key}' value='#{value}'/><br />"
    when 'break'
      html_return = '<br/>'
    when 'list'
      html_return = "#{data['description']}: <select name='#{key}'>"
      value = data['value'].to_s empty? data['default_value'] : data['value']
      data['data'].each do |list|
        selected = value == list['value'] ? 'selected' : ''
        html_return += "<option value='#{list['value']}' #{selected}>#{list['text']}</option>"
      end
      html_return += '</select><br />'
    when 'radio'
      html_return = "#{data['description']}:"
      data['data'].each do |list|
        checked = list['checked'] ? 'checked' : ''
        html_return += "|<input type='radio' name='#{key}' value='#{key}' #{checked}/#{list['description']}"
      end
      html_return += '<br />'
    when 'checkbox'
      value = data['value'].to_s.empty? ? data['default_value'] : data['value']
      checked = value ? 'checked' : ''
      html_return = "#{data['description']}: <input type='checkbox' name='#{key}' #{checked}/><br />"
    end
    html_return
  end

  # Reads a JSON file and returns its contents as a Hash.
  def file2json(file)
    raise 'cant find file' unless File.exist?(file)

    data = File.read(file)
    JSON.parse(data)
  end
end
