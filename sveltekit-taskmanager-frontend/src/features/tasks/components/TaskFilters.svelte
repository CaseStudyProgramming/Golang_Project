<script lang="ts">
  const { 
    searchQuery, 
    filterCompleted 
  }: { 
    searchQuery: { value: string }, 
    filterCompleted: { value: boolean | undefined } 
  } = $props();
</script>

<div class="bg-white rounded-lg shadow p-4">
  <div class="flex flex-col sm:flex-row gap-4">
    <div class="flex-1">
      <input 
        type="text" 
        placeholder="Search tasks..." 
        bind:value={searchQuery.value}
        class="input"
      />
    </div>
    <div class="flex gap-2">
      <button 
        class="btn btn-secondary"
        class:btn-primary={filterCompleted.value === undefined}
        onclick={() => filterCompleted.value = undefined}
      >
        All
      </button>
      <button 
        class="btn btn-secondary"
        class:btn-primary={filterCompleted.value === false}
        onclick={() => filterCompleted.value = false}
      >
        Active
      </button>
      <button 
        class="btn btn-secondary"
        class:btn-primary={filterCompleted.value === true}
        onclick={() => filterCompleted.value = true}
      >
        Completed
      </button>
    </div>
  </div>
</div>
