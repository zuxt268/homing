package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/dto/req"
	"github.com/zuxt268/homing/internal/interface/dto/res"
	_ "github.com/zuxt268/homing/internal/interface/dto/res"
	"github.com/zuxt268/homing/internal/usecase"
)

type APIHandler struct {
	customerUsecase           usecase.CustomerUsecase
	tokenUsecase              usecase.TokenUsecase
	wordpressInstagramUsecase usecase.WordpressInstagramUsecase
	businessInstagramUsecase  usecase.BusinessInstagramUsecase
}

func NewAPIHandler(
	customerUsecase usecase.CustomerUsecase,
	tokenUsecase usecase.TokenUsecase,
	wordpressInstagramUsecase usecase.WordpressInstagramUsecase,
	businessInstagramUsecase usecase.BusinessInstagramUsecase,
) APIHandler {
	return APIHandler{
		customerUsecase:           customerUsecase,
		tokenUsecase:              tokenUsecase,
		wordpressInstagramUsecase: wordpressInstagramUsecase,
		businessInstagramUsecase:  businessInstagramUsecase,
	}
}

// SyncAllGoogleBusinessInstagram godoc
// @Summary      instagram => wordpressにおける全顧客データ同期
// @Description  全ての顧客のデータを同期します
// @Tags         sync
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "全顧客同期完了"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/sync/business-instagram [post]
func (h *APIHandler) SyncAllGoogleBusinessInstagram(c echo.Context) error {
	fmt.Println("aaaa")
	err := h.customerUsecase.SyncAllGoogleBusinessInstagram(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "sync all")
}

// SyncOneGoogleBusinessInstagram godoc
// @Summary      instagram => wordpressにおける顧客データ同期
// @Description  全ての顧客のデータを同期します
// @Tags         sync
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "全顧客同期完了"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/sync/business-instagram/{id} [post]
func (h *APIHandler) SyncOneGoogleBusinessInstagram(c echo.Context) error {
	var id int
	if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err := h.customerUsecase.SyncOneGoogleBusinessInstagram(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "sync one")
}

// SyncAllWordpressInstagram godoc
// @Summary      instagram => wordpressにおける全顧客データ同期
// @Description  全ての顧客のデータを同期します
// @Tags         sync
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "全顧客同期完了"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/sync/wordpress-instagram [post]
func (h *APIHandler) SyncAllWordpressInstagram(c echo.Context) error {
	err := h.customerUsecase.SyncAllWordpressInstagram(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "sync all")
}

// SyncOneWordpressInstagram godoc
// @Summary      instagram => wordpressにおける顧客データ同期
// @Description  全ての顧客のデータを同期します
// @Tags         sync
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "全顧客同期完了"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/sync/wordpress-instagram/{id} [post]
func (h *APIHandler) SyncOneWordpressInstagram(c echo.Context) error {
	var id int
	if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err := h.customerUsecase.SyncOneWordpressInstagram(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "sync one")
}

// SaveToken godoc
// @Summary      トークンを保存します。
// @Description
// @Tags         token
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "更新完了"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/token [post]
func (h *APIHandler) SaveToken(c echo.Context) error {
	var token req.UpdateToken
	if err := c.Bind(&token); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := h.tokenUsecase.UpdateToken(c.Request().Context(), token)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "ok")
}

// GetToken godoc
// @Summary      トークンを取得します。
// @Description
// @Tags         token
// @Accept       json
// @Produce      json
// @Success　　　 200 {object} res.Token "トークン情報"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/token [get]
func (h *APIHandler) GetToken(c echo.Context) error {
	token, err := h.tokenUsecase.GetToken(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, token)
}

// CheckToken godoc
// @Summary      トークンの認証情報を取得する
// @Description
// @Tags         token
// @Accept       json
// @Produce      json
// @Success　　　 200 {object} string "ok"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/token/check [post]
func (h *APIHandler) CheckToken(c echo.Context) error {
	err := h.tokenUsecase.CheckToken(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "ok")
}

// GetWordpressInstagramList godoc
// @Summary      Wordpress Instagram一覧取得
// @Description  Wordpress Instagramの一覧を取得します
// @Tags         wordpress-instagram
// @Accept       json
// @Produce      json
// @Param        limit          query     int     false  "取得件数"
// @Param        offset         query     int     false  "オフセット"
// @Param        name           query     string  false  "名前"
// @Param        wordpress      query     string  false  "WordPress URL"
// @Param        instagram_id   query     string  false  "Instagram ID"
// @Param        status         query     int     false  "ステータス"
// @Param        delete_hash    query     bool    false  "削除ハッシュ"
// @Param        customer_type  query     int     false  "顧客タイプ"
// @Success      200  {object}  res.WordpressInstagramList  "Wordpress Instagram一覧"
// @Failure      400  {string}  string  "不正なリクエスト"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/wordpress-instagram [get]
func (h *APIHandler) GetWordpressInstagramList(c echo.Context) error {
	var params req.GetWordpressInstagram
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	list, err := h.wordpressInstagramUsecase.GetWordpressInstagramList(c.Request().Context(), params)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, list)
}

// GetWordpressInstagramCount godoc
// @Summary      Wordpress Instagramの件数を取得
// @Description  Wordpress Instagramの件数を取得します
// @Tags         wordpress-instagram
// @Accept       json
// @Produce      json
// @Success      200  {object}  res.WordpressInstagramList  "Wordpress Instagram一覧"
// @Failure      400  {string}  string  "不正なリクエスト"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/wordpress-instagram/count [get]
func (h *APIHandler) GetWordpressInstagramCount(c echo.Context) error {
	list, err := h.wordpressInstagramUsecase.GetWordpressInstagramCount(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, list)
}

// GetWordpressInstagram godoc
// @Summary      Wordpress Instagram詳細取得
// @Description  Wordpress Instagramの詳細を取得します
// @Tags         wordpress-instagram
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Wordpress Instagram ID"
// @Param        limit          query     int     false  "投稿取得件数"
// @Param        offset         query     int     false  "投稿オフセット"
// @Success      200  {object}  res.WordpressInstagramDetail  "Wordpress Instagram詳細"
// @Failure      404  {string}  string  "見つかりません"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/wordpress-instagram/{id} [get]
func (h *APIHandler) GetWordpressInstagram(c echo.Context) error {
	var params req.GetWordpressInstagramDetail
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	var id int
	if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	item, err := h.wordpressInstagramUsecase.GetWordpressInstagram(c.Request().Context(), id, params)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, item)
}

// CreateWordpressInstagram godoc
// @Summary      Wordpress Instagram作成
// @Description  Wordpress Instagramを作成します
// @Tags         wordpress-instagram
// @Accept       json
// @Produce      json
// @Param        body  body      req.CreateWordpressInstagram  true  "作成データ"
// @Success      201   {object}  res.WordpressInstagram  "作成されたWordpress Instagram"
// @Failure      400   {string}  string  "不正なリクエスト"
// @Failure      500   {string}  string  "内部サーバーエラー"
// @Router       /api/wordpress-instagram [post]
func (h *APIHandler) CreateWordpressInstagram(c echo.Context) error {
	var body req.CreateWordpressInstagram
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	item, err := h.wordpressInstagramUsecase.CreateWordpressInstagram(c.Request().Context(), body)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusCreated, item)
}

// UpdateWordpressInstagram godoc
// @Summary      Wordpress Instagram更新
// @Description  Wordpress Instagramを更新します
// @Tags         wordpress-instagram
// @Accept       json
// @Produce      json
// @Param        id    path      int                           true  "Wordpress Instagram ID"
// @Param        body  body      req.UpdateWordpressInstagram  true  "更新データ"
// @Success      200   {object}  res.WordpressInstagram  "更新されたWordpress Instagram"
// @Failure      400   {string}  string  "不正なリクエスト"
// @Failure      404   {string}  string  "見つかりません"
// @Failure      500   {string}  string  "内部サーバーエラー"
// @Router       /api/wordpress-instagram/{id} [put]
func (h *APIHandler) UpdateWordpressInstagram(c echo.Context) error {
	var body req.UpdateWordpressInstagram
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var id int
	if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	item, err := h.wordpressInstagramUsecase.UpdateWordpressInstagram(c.Request().Context(), id, body)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, item)
}

// DeleteWordpressInstagram godoc
// @Summary      Wordpress Instagram削除
// @Description  Wordpress Instagramを削除します
// @Tags         wordpress-instagram
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Wordpress Instagram ID"
// @Success      204  {string}  string  "削除成功"
// @Failure      404  {string}  string  "見つかりません"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/wordpress-instagram/{id} [delete]
func (h *APIHandler) DeleteWordpressInstagram(c echo.Context) error {
	var id int
	if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := h.wordpressInstagramUsecase.DeleteWordpressInstagram(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// FetchGoogleBusinessList godoc
// @Summary      Google Businessの同期
// @Description  Google Businessを同期します
// @Tags         google-business
// @Failure      400  {string}  string  "不正なリクエスト"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/google-business/fetch [post]
func (h *APIHandler) FetchGoogleBusinessList(c echo.Context) error {
	err := h.businessInstagramUsecase.FetchGoogleBusinesses(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "ok")
}

// GetGoogleBusinessList godoc
// @Summary      Google Business一覧取得
// @Description  Google Businessの一覧を取得します（ページング対応）
// @Tags         google-business
// @Accept       json
// @Produce      json
// @Param        limit   query     int  false  "取得件数（デフォルト: 20）"
// @Param        offset  query     int  false  "オフセット（デフォルト: 0）"
// @Success      200  {object}  res.GoogleBusinessList  "Google Business一覧"
// @Failure      400  {string}  string  "不正なリクエスト"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/google-business [get]
func (h *APIHandler) GetGoogleBusinessList(c echo.Context) error {
	limit := 20
	offset := 0

	if err := echo.QueryParamsBinder(c).
		Int("limit", &limit).
		Int("offset", &offset).
		BindError(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	businesses, total, err := h.businessInstagramUsecase.GetGoogleBusinesses(c.Request().Context(), limit, offset)
	if err != nil {
		return handleError(c, err)
	}

	businessList := make([]res.GoogleBusiness, 0, len(businesses))
	for _, b := range businesses {
		businessList = append(businessList, res.GoogleBusiness{
			ID:        b.ID,
			Name:      b.Name,
			Title:     b.Title,
			CreatedAt: b.CreatedAt,
		})
	}

	return c.JSON(http.StatusOK, res.GoogleBusinessList{
		GoogleBusinessList: businessList,
		Paginate: res.Paginate{
			Total: total,
			Count: len(businesses),
		},
	})
}

// GetBusinessInstagramList godoc
// @Summary      Business Instagram一覧取得
// @Description  Business Instagram一覧を取得します
// @Tags         business-instagram
// @Accept       json
// @Produce      json
// @Success      201   {object}  res.BusinessInstagramList  "Business Instagram"
// @Failure      400   {string}  string  "不正なリクエスト"
// @Failure      500   {string}  string  "内部サーバーエラー"
// @Router       /api/business-instagram [get]
func (h *APIHandler) GetBusinessInstagramList(c echo.Context) error {
	var params req.GetBusinessInstagram
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	resp, err := h.businessInstagramUsecase.GetBusinessInstagramList(c.Request().Context(), params)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusCreated, resp)
}

// GetBusinessInstagram godoc
// @Summary      Business Instagram取得
// @Description  Business Instagramを取得します
// @Tags         business-instagram
// @Accept       json
// @Produce      json
// @Success      201   {object}  res.BusinessInstagram  "Business Instagram"
// @Failure      400   {string}  string  "不正なリクエスト"
// @Failure      500   {string}  string  "内部サーバーエラー"
// @Router       /api/business-instagram/{id} [get]
func (h *APIHandler) GetBusinessInstagram(c echo.Context) error {
	var id int
	if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	resp, err := h.businessInstagramUsecase.GetBusinessInstagram(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

// CreateBusinessInstagram godoc
// @Summary      Business Instagram作成
// @Description  Business Instagramを作成します
// @Tags         business-instagram
// @Accept       json
// @Produce      json
// @Param        body  body      req.BusinessInstagram  true  "作成データ"
// @Success      201   {object}  res.BusinessInstagram  "作成されたBusiness Instagram"
// @Failure      400   {string}  string  "不正なリクエスト"
// @Failure      500   {string}  string  "内部サーバーエラー"
// @Router       /api/business-instagram [post]
func (h *APIHandler) CreateBusinessInstagram(c echo.Context) error {
	var body req.BusinessInstagram
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	fmt.Println(body)

	resp, err := h.businessInstagramUsecase.CreateBusinessInstagram(c.Request().Context(), body)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusCreated, resp)
}

// UpdateBusinessInstagram godoc
// @Summary      Business Instagram更新
// @Description  Business Instagramを更新します
// @Tags         business-instagram
// @Accept       json
// @Produce      json
// @Param        id    path      int                    true  "Business Instagram ID"
// @Param        body  body      req.BusinessInstagram  true  "更新データ"
// @Success      200   {object}  res.BusinessInstagram  "更新されたBusiness Instagram"
// @Failure      400   {string}  string  "不正なリクエスト"
// @Failure      404   {string}  string  "見つかりません"
// @Failure      500   {string}  string  "内部サーバーエラー"
// @Router       /api/business-instagram/{id} [put]
func (h *APIHandler) UpdateBusinessInstagram(c echo.Context) error {
	var id int
	if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	var body req.BusinessInstagram
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	resp, err := h.businessInstagramUsecase.UpdateBusinessInstagram(c.Request().Context(), id, body)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

// DeleteBusinessInstagram godoc
// @Summary      Business Instagram削除
// @Description  Business Instagramを削除します
// @Tags         business-instagram
// @Param        id    path      int     true  "Business Instagram ID"
// @SuccessWI      204
// @Failure      400   {string}  string  "不正なリクエスト"
// @Failure      404   {string}  string  "見つかりません"
// @Failure      500   {string}  string  "内部サーバーエラー"
// @Router       /api/business-instagram/{id} [delete]
func (h *APIHandler) DeleteBusinessInstagram(c echo.Context) error {
	var id int
	if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := h.businessInstagramUsecase.DeleteBusinessInstagram(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func handleError(c echo.Context, err error) error {
	slog.Error("handleError", "error", err.Error())
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.JSON(http.StatusNotFound, res.ErrorResponse{Message: err.Error()})
	case errors.Is(err, domain.ErrBadRequest):
		return c.JSON(http.StatusBadRequest, res.ErrorResponse{Message: err.Error()})
	case errors.Is(err, domain.ErrWordpressConnection):
		return c.JSON(http.StatusBadRequest, res.ErrorResponse{Message: err.Error()})
	case errors.Is(err, domain.ErrInstagramConnection):
		return c.JSON(http.StatusBadRequest, res.ErrorResponse{Message: err.Error()})
	default:
		return c.JSON(http.StatusInternalServerError, res.ErrorResponse{Message: err.Error()})
	}
}
