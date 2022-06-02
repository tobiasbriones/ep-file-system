<template>
  <Home />
</template>

<script>
import Home from './components/Home.vue';

export default {
  name: 'App',
  components: {
    Home: Home
  }
};
</script>

<style>
#app {
  margin-top: 60px;
  font-family: Avenir, Helvetica, Arial, sans-serif;
  text-align: center;
  color: #2c3e50;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
