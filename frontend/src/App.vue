<template>
  <div id="app">
    <UrlShortener/>
  </div>
</template>

<script>
  import UrlShortener from "./components/UrlShortener";

  export default {
    name: "App",
    components: {
      UrlShortener
    }
  };
</script>

<style>
  #app {
    font-family: "Avenir", Helvetica, Arial, sans-serif;
    font-size: 18px;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: left;
    color: #2c3e50;
    margin: 60px;
    padding: 10px;
  }
</style>
