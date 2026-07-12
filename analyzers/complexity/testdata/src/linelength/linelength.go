package linelength

func good() { _ = 1 }

func goodString() { _ = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" }

func bad() { _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; _ = 1; } // want "line is .* characters long"
