<script setup>
import {ref} from 'vue'
import SvgIcon from "@/components/SvgIcon.vue";

const emit = defineEmits(['confirm', 'cancel'])
const props = defineProps(['question'])

const image = ref(null)
const caption = ref('')

const onConfirm = () => {
    emit('confirm', {image: image.value, mimeType: image.value.type, caption: caption.value})
    image.value = null
    caption.value = ''
}

const onCancel = () => {
    emit('cancel')
    image.value = null
    caption.value = ''
}

const onFileChange = (event) => {
    image.value = event.target.files[0]
}
</script>

<template>
    <div class="modal-backdrop">
        <div class="modal-container">
            <form @submit.prevent="onConfirm" @reset.prevent="onCancel">
                <div class="modal-body">
                    <h2>{{ question }}</h2>
                    <input type="file"
                           class="form-control form-control-sm"
                           accept="image/*"
                           @change="onFileChange"
                           required>
                    <input type="text"
                           class="form-control form-control-sm"
                           placeholder="Caption (optional)"
                           v-model.trim="caption"
                           pattern="^[\p{L}\p{N}\p{M}\p{P}\p{S} ]{0,32}$"
                           minlength="0"
                           maxlength="32"
                           title="0 to 32 characters (UNICODE supported)">
                    <div class="modal-action">
                        <button type="submit" class="modal-button btn btn-sm btn-outline-success"
                                :disabled="!image || caption.trim().length > 32">Confirm
                            <svg-icon icon="check-square"/>
                        </button>
                        <button type="reset" class="modal-button btn btn-sm btn-outline-danger">Cancel
                            <svg-icon icon="x-square"/>
                        </button>
                    </div>
                </div>
            </form>
        </div>
    </div>
</template>

<style scoped>
.modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
}

.modal-container {
    background-color: white;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
    z-index: 1000;
    min-width: max-content;
}

.modal-body {
    display: flex;
    flex-direction: column;
    gap: 5px;
}

.modal-action .modal-button:not(:last-child) {
    margin-right: 10px;
}
</style>
