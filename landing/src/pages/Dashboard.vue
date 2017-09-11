<template>
  <div id="app">
    <Header></Header>
    <Chart :item="item"></Chart>
    <Search :onItemSelect="onItemSelect"></Search>
    <Footer></Footer>
  </div>
</template>

<script>
  import Header from '../components/Header'
  import Chart from '../components/Chart'
  import Search from '../components/Search'
  import Footer from '../components/Footer'

  export default {
    name: 'dashboard',
    data(){
      return {item: null,};
    },
    components: {
      Header,
      Chart,
      Search,
      Footer,
    },
    methods:{
      onItemSelect(item){
        console.log('item select')
        this.item = item;
      }
    }
  }
</script>

<style>
  @import '../style/index.css';
</style>
