<script>
  import { nations, tiers, shipClassImages } from './wows';
  export let nation;
  export let type;
  export let tier;
  export let name;
</script>

<div class="flex text-center">
  <img alt={`${type} icon`} src={shipClassImages[type]} />
  {tiers[tier - 1]} {name}
</div>
