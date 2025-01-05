using UnityEngine;

public class BarrierController : MonoBehaviour
{
    public float speed = 2f;
    public float lifetime = 10f;

    private void Start()
    {
        //实现超出范围后，销毁多余的障碍物
        Destroy(gameObject, lifetime);
    }

    private void Update()
    {
        MoveLeft();
    }

    private void MoveLeft()
    {
        //向左移动障碍物
        // Debug.Log("move left....");
        transform.Translate(Vector2.left * speed * Time.deltaTime);
    }
}
