<script>
  export let name = "";
  export let placeholder = "";
  export let type = "text";
</script>

{#if type !== "textarea"}
  <input
    class="p-4 border border-green-500 rounded"
    type={type}
    placeholder={placeholder}
    name={name}
  />
{:else}
  <textarea
    class="p-4 border border-green-500 rounded"
    type={type}
    placeholder={placeholder}
    name={name}
  ></textarea>
{/if}
