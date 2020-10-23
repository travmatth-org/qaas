<template>
  <div class='full-width center-content'>
    <hello-component name='World' />
  </div>
</template>

<script lang='ts'>
import HelloComponent from '../components/Hello.vue'

export default {
  components: {
    HelloComponent
  }
};
</script>

<style lang='scss'>
.full-width {
  width: 100%;
}

.center-content {
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>