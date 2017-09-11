<template>
  <div id="app">
    <Header></Header>
    <Hero></Hero>
    <Items></Items>
    <Timeline></Timeline>
    <Footer></Footer>
  </div>
</template>

<script>
import Header from '../components/Header'
import Hero from '../components/Hero'
import Items from '../components/Items'
import Timeline from '../components/Timeline'
import Footer from '../components/Footer'

export default {
  name: 'app',
  components: {
    Header,
    Hero,
    Items,
    Timeline,
    Footer,
  }
}
</script>

<style>
@import '../style/index.css';
</style>
