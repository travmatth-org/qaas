<template>
  <h1>Hello {{ name }}!</h1>
</template>

<script lang='ts'>
export default {
  props: {
    name: {
      type: String,
      required: true
    }
  }
};
</script>

<style lang='scss'>
h1 {
  color: $primary;
}
</style>