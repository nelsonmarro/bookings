function changeProcessed(src, id, year, month) {
  window.location.href = `/admin/reservations/${src}/${id}/process?y=${year}&m=${month}`;
}

function deleteReservation(src, id, year, month) {
  window.location.href = `/admin/reservations/${src}/${id}/delete?y=${year}&m=${month}`;
}
