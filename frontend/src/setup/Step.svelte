<script>
  export let number;
  export let title;

  export let active = false;
  export let done = false;
</script>

<style>
.step-circle {
  transition: border-color .2s, color .2s;
}
</style>

<div class="flex items-center mb-4" class:text-gray-500={active === false}>
  <div class="step-circle border-2 px-4 py-2 rounded-full font-bold text-blue mr-2" class:text-green-400={done !== false} class:border-gray-900={active} class:border-green-400={done !== false}>{number}</div>
  <h2 class="text-lg">{title}</h2>
</div>