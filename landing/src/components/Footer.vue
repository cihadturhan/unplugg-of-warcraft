<template>
  <footer>
    <div class="container">
      <div class="row justify-content-center color-light">
        <div class="col-md-2"><img src="../assets/logo.svg" alt=""></div>
        <div class="col-md-2">
          <h3>Company</h3>
          <ul>
            <li>Contact</li>
            <li>About Us</li>
            <li>Privacy Policy</li>
            <li>Pricing</li>
          </ul>
        </div>
        <div class="col-md-2">
          <h3>Social Media</h3>
          <ul>
            <li>Twitter</li>
            <li>Facebook</li>
          </ul>
        </div>
      </div>
    </div>
  </footer>
</template>

<style>
  footer{
    background-color: var(--color-dark-blue-gray);
    background-image: linear-gradient(hsla(0, 0%, 0%, 0.1), hsla(0, 0%, 0%, 0.4));
    min-height: 100px;
    position: relative;
    overflow: hidden;
    padding: 20px 0;
  }

  footer>.container{
    position: relative;
    z-index: 1;
  }

  footer::before{
    content: "";
    position: absolute;
    box-shadow: 0 1px 2px black, 0 0px 50px hsla(0, 0%, 0%, 0.9);
    height: 50px;
    left: 0;
    top: -50px;
    width: 100%;
    animation: slowUpDown 1.8s infinite both alternate ease;
  }

  li{
    font-size: 18px;
    font-weight: 300;
    line-height: 1.4;
  }

  @keyframes slowUpDown{
    from{
      transform: translateY(-50px)
    }
    to{
      transform: translateY(0px)
    }
  }
</style>
