<template>
  <div class="container">
    <div class="row justify-content-center ">
      <div class="search-container">
        <input type="text" @focus="showDropdown()">
        <div class="search-dropdown" v-show="isDropdownShown">
          <ul>
            <li @click="_onItemSelect(item)" v-for="item in items"> {{item.name}}</li>
          </ul>
        </div>
        <button class="sm">Search</button>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    props: ['onItemSelect'],
    data(){
      return {
        isDropdownShown: false,
        items: [
          { name: 'Helm of the Demonic Gladiator',},
          { name: 'Legion Season 7 EliteChest of the Demonic Gladiator',},
          { name: 'Leggings of the Demonic Gladiator',},
          { name: 'Pauldrons of the Demonic Gladiator',},
          { name: 'Bracers of the Dominating Gladiator',},
          { name: 'Helm of the Awakened Soul',},
          ]
      }
    },
    methods: {
      _onItemSelect(item){
        this.onItemSelect(item);
        this.isDropdownShown = false;
      },
      showDropdown(){
        this.isDropdownShown = true;
      },
    }
  }
</script>

<style>
  .search-container{
    padding: 30px 0;
    display: flex;
    align-items: center;
    position: relative;
  }

  .search-dropdown{
    position: absolute;
    top: 80px;
    width: 100%;
    background: white;
    z-index: 10;
  }

  .search-dropdown li{
    cursor: pointer;
    padding: 10px;
  }

  .search-dropdown li:hover{
    background: hsl(0, 0%, 96%);
  }

  input{
    appearance: none;
    border: 1px solid var(--color-red);
    min-height: 40px;
    border-radius: 30px;
  }
</style>
