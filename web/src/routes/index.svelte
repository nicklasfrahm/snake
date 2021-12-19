<script>
  import Input from '../components/Input.svelte'
</script>

<div class="w-screen h-screen bg-gray-50 flex justify-center items-center">
  <div class="max-w-md p-8 bg-white border border-green-500 rounded-xl grid gap-4">
    <h1 class="font-bold text-green-500 text-center text-3xl">Snek</h1>
    <h2 class="font-medium text-center text-xl">Publish your queues online</h2>

    <p class="text-justify">Are you hosting an event and the queue is huge?</p>
    <p class="text-justify">People are freezing outside?</p>

    <p class="text-justify">Then this free tool is right for you! Simply fill out the form below to create a queue and publish it online. Everyone can see the queue number online, which allows them to do other things while waiting.</p>

    <Input placeholder="Organizer email" name="owner" />
    <Input placeholder="Event title" name="title" />
    <Input placeholder="Description" name="description" type="textarea" />
    <Input placeholder="Start number" name="number" />

    <button class="px-4 py-2 transition-color bg-green-500 hover:bg-green-600 focus:bg-green-600 active:bg-green-500 text-white font-medium rounded uppercase">
      Create queue
    </button>
  </div>
</div>
