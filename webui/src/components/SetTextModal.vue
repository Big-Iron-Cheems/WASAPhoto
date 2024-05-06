<script setup>
import {ref} from 'vue'
import SvgIcon from "@/components/SvgIcon.vue";

const emit = defineEmits(['confirm', 'cancel'])
const props = defineProps({
    'header': String,
    'placeholder': String,
    'pattern': String,
    'minlength': Number,
    'maxlength': Number,
    'title': String,
    'rows': Number,
    'allowEnter': Boolean,
})

const inputText = ref('')
const regex = new RegExp(props.pattern, "u")

const handleEnterKey = (event) => {
    if (!props.allowEnter) event.preventDefault()
}
const validateInput = (inputText) => {
    return inputText.length >= props.minlength
        && inputText.length <= props.maxlength
        && regex.test(inputText)
}
const onConfirm = () => {
    if (validateInput(inputText.value)) {
        emit('confirm', {inputText: inputText.value})
        inputText.value = ''
    }
}
const onCancel = () => {
    emit('cancel')
    inputText.value = ''
}
</script>

<template>
    <div class="modal-backdrop">
        <div class="modal-container">
            <form @submit.prevent="onConfirm" @reset.prevent="onCancel">
                <div class="modal-body">
                    <h2>{{ $props.header }}</h2>
                    <textarea type="text"
                              class="form-control form-control-sm"
                              :placeholder="$props.placeholder"
                              v-model.trim="inputText"
                              required
                              :minlength="$props.minlength"
                              :maxlength="$props.maxlength"
                              :title="$props.title"
                              :rows="$props.rows"
                              cols="32"
                              v-on:keydown.enter="handleEnterKey"
                    />
                    <div class="modal-action">
                        <button type="submit"
                                class="modal-button btn btn-sm btn-outline-success"
                                :disabled="!validateInput(inputText)">
                            Confirm
                            <svg-icon icon="check-square"/>
                        </button>
                        <button type="reset"
                                class="modal-button btn btn-sm btn-outline-danger">
                            Cancel
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
