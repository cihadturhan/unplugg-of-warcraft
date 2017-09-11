<template>
  <section class="items">
    <div class="container">
      <h2>SEE FUTURE PRICE OF ANY ITEMS NOW</h2>
      <div class="row justify-content-center">
        <div class="col-md-10">
          <p> Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed rhoncus erat vitae felis luctus commodo. Morbi mattis aliquam turpis sed porta. Nullam vel imperdiet massa. Curabitur blandit ultricies volutpat. Morbi sit amet massa ut ligula rhoncus pharetra id in massa. Quisque et nunc elit. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas.</p>
        </div>
      </div>
      <div class="row justify-content-center">
        <div class="col-md-2">
          <div class="item"></div>
        </div>

        <div class="col-md-2">
          <div class="item"></div>
        </div>

        <div class="col-md-2">
          <div class="item"></div>
        </div>

        <div class="col-md-2">
          <div class="item"></div>
        </div>

        <div class="col-md-2">
          <div class="item"></div>
        </div>


      </div>

      <div class="row justify-content-center mar-top-md">
        <a class="btn btn-primary">Get Your Items</a>
      </div>
    </div>
  </section>
</template>

<style>
  .items{
    width: 100vw;
    height: 100vh;
    text-align: center;
    background-image: url('../assets/image-2.png');
    display: flex;
    align-items: center;
  }

  .item{
    height: 80px;
    width: 100%;
    border-radius: 6px;
    background: var(--color-light);
    box-shadow:  0 1px 2px hsla(0, 0%, 0%, 0.2), 0 6px 18px hsla(0, 0%, 0%, 0.2);
  }

  .mar-top-md{
    margin-top: 60px;
  }

</style>
