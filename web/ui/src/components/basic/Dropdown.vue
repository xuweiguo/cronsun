<template>
  <div class="cs-dropdown">
    <select
      class="layui-input"
      :multiple="isMultiple"
      v-model="internalValue"
      @change="onChange">
      <option
        v-if="placeholderText && !isMultiple"
        disabled
        value="">
        {{placeholderText}}
      </option>
      <option
        v-for="option in optionList"
        :key="option.value"
        :value="option.value">
        {{option.name}}
      </option>
    </select>
    <input
      v-if="allowAdditions"
      type="text"
      class="layui-input dropdown-addition"
      :placeholder="$L('add and press enter')"
      @keyup.enter="addCustom">
  </div>
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
      internalValue: this.getInitialValue(),
      additions: []
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

    optionList() {
      return [...this.normalizedItems, ...this.additions];
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

    onChange() {
      this.$emit('change', this.internalValue);
    },

    addCustom(event) {
      const value = (event.target.value || '').trim();
      if (!value) return;
      if (!this.optionList.find((opt) => opt.value === value)) {
        this.additions.push({ value, name: value });
      }
      if (this.isMultiple) {
        const values = Array.isArray(this.internalValue) ? [...this.internalValue] : [];
        if (!values.includes(value)) {
          values.push(value);
        }
        this.internalValue = values;
      } else {
        this.internalValue = value;
      }
      event.target.value = '';
      this.onChange();
    }
  }
}
</script>

<style scoped>
.cs-dropdown {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.dropdown-addition {
  font-size: 12px;
}
</style>
