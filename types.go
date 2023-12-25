package main

type GelbooruPosts struct {
	Attributes struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Count  int `json:"count"`
	} `json:"@attributes"`
	Post []struct {
		ID            int    `json:"id"`
		CreatedAt     string `json:"created_at"`
		Score         int    `json:"score"`
		Width         int    `json:"width"`
		Height        int    `json:"height"`
		Md5           string `json:"md5"`
		Directory     string `json:"directory"`
		Image         string `json:"image"`
		Rating        string `json:"rating"`
		Source        string `json:"source"`
		Change        int    `json:"change"`
		Owner         string `json:"owner"`
		CreatorID     int    `json:"creator_id"`
		ParentID      int    `json:"parent_id"`
		Sample        int    `json:"sample"`
		PreviewHeight int    `json:"preview_height"`
		PreviewWidth  int    `json:"preview_width"`
		Tags          string `json:"tags"`
		Title         string `json:"title"`
		HasNotes      string `json:"has_notes"`
		HasComments   string `json:"has_comments"`
		FileURL       string `json:"file_url"`
		PreviewURL    string `json:"preview_url"`
		SampleURL     string `json:"sample_url"`
		SampleHeight  int    `json:"sample_height"`
		SampleWidth   int    `json:"sample_width"`
		Status        string `json:"status"`
		PostLocked    int    `json:"post_locked"`
		HasChildren   string `json:"has_children"`
	} `json:"post"`
}
