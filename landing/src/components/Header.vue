<template>
  <header>
    <img src="../assets/logo.svg" alt="">
    <nav>
      <a href="#">Home</a>
      <a href="#">Items</a>
      <a href="#">Contact</a>
      <a href="#">Login</a>

      <a href="#" class="primary">Get Started</a>
    </nav>
  </header>
</template>


<style>
  header{
    position: fixed;
    left: 0;
    top: 0;
    right: 0;
    padding: 15px;
    background: hsla(0, 0%, 0%, 0.2);
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  nav>a{
    color: var(--color-light);
    text-decoration: none;
    font-size: 22px;
    margin-left: 20px;
    opacity: 0.85;
  }

  nav>a.primary{
    font-weight: bold;
    opacity: 0.9;
    margin-left: 30px;
  }

  nav>a:hover{
    opacity: 1;
  }
</style>
