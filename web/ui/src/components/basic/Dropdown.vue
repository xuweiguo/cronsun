<template>
  <el-select
    class="cs-dropdown"
    :multiple="isMultiple"
    :placeholder="placeholderText"
    :value="internalValue"
    :filterable="true"
    :allow-create="!!allowAdditions"
    default-first-option
    @input="onChange"
    style="width: 100%">
    <el-option
      v-for="option in normalizedItems"
      :key="option.value"
      :label="option.name"
      :value="option.value"/>
  </el-select>
</template>

<script>
export default {
  name: 'dropdown',
  props: {
    title: { type: String, default: '' },
    items: { type: Array, default: () => [] },
    allowAdditions: { type: [Boolean, String], default: false },
    multiple: { type: [Boolean, String], default: false },
    selected: { type: [String, Number, Array], default: null }
  },

  data() {
    return {
      internalValue: this.getInitialValue()
    }
  },

  computed: {
    isMultiple() {
      return this.multiple === '' ? true : !!this.multiple;
    },

    normalizedItems() {
      return (this.items || []).map((item) => {
        if (typeof item === 'object') {
          return { value: item.value, name: item.name };
        }
        return { value: item, name: item };
      });
    },

    placeholderText() {
      if (this.title && this.title.length > 0) {
        return this.title;
      }
      return this.isMultiple ? this.$L('select groups') : this.$L('select a group');
    }
  },

  watch: {
    selected: {
      handler(val) {
        this.internalValue = val || (this.isMultiple ? [] : '');
      },
      deep: true
    }
  },

  methods: {
    getInitialValue() {
      return this.selected || (this.isMultiple ? [] : '');
    },

    onChange(value) {
      this.internalValue = value;
      this.$emit('change', value);
    }
  }
}
</script>
