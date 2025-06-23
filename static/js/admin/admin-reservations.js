// variables

// functions
function loadTable(data, reservationType) {
  new gridjs.Grid({
    columns: [
      {
        id: "ID",
        name: "ID",
      },
      {
        id: "FirstName",
        name: "First Name",
        formatter: (_, row) =>
          gridjs.html(
            `<a class="hover:underline text-blue-500" href="/admin/reservations/${reservationType}/${row.cells[0].data}">${row.cells[1].data}</a>`,
          ),
      },
      {
        id: "LastName",
        name: "Last Name",
      },
      {
        data: (row) => row.Room.RoomName,
        name: "Room",
      },
      {
        data: (row) => new Date(row.StartDate).toLocaleDateString(),
        name: "Arrival",
      },
      {
        data: (row) => new Date(row.EndDate).toLocaleDateString(),
        name: "Departure",
      },
    ],
    data: data,
    search: true,
    sort: true,
    pagination: {
      limit: 5,
    },
  }).render(document.getElementById("res-table"));
}
